# Graph Report - app-listener  (2026-10-12)

## Corpus Check
- 347 files · ~1,550,945 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 12 file(s) not represented in the graph (top: (none) 7, .service 2, .pub 1)

## Summary
- 4560 nodes · 17962 edges · 212 communities (170 shown, 42 thin omitted)
- Extraction: 93% EXTRACTED · 7% INFERRED · 0% AMBIGUOUS · INFERRED: 1301 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `989e5ddd`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- daemon/daemon.go
- testing.T
- IntegrationSuite
- filevault.go
- daemonconfig_test.go
- networkguard.bpf.c
- monitor.bpf.c
- User
- guardUnitTest
- mkdirs
- Hash
- planUpdaters
- codeedit.go
- IntegrationSuite
- RunCmd
- Vault
- IntegrationSuite
- guard_trust.bpf.c
- RefreshSection
- fileEditModel
- networkmonitor.bpf.c
- usecase/daemon_test.go
- Resource
- RawTarget
- runNetworkGuard
- guard.bpf.c
- NetworkMonitor
- monitorUnitTest
- runDaemon
- buildOneGuard
- Monitor
- serve.go
- engine
- IntegrationSuite
- fileedit_model_test.go
- daemonModel
- watchGroup
- filevault_test.go
- GuardEvent
- guardSpec
- github.com/charmbracelet/bubbles/textarea.Model
- runNetworkMonitor
- absPath
- liveSession
- IntegrationSuite
- SelectOrphans
- IntegrationSuite
- .syncSetsLocked
- TrustGuard
- tempRow
- __always_inline
- IntegrationSuite
- MountVouch
- github.com/charmbracelet/bubbletea.Cmd
- IntegrationSuite
- highlight.go
- sync.Mutex
- GenerateConf
- EventType
- fileedit_symlink_test.go
- ListUsers
- exe_supersede.h
- fakeGuardRepo
- errno
- pinstate.go
- controlServer
- time.Duration
- LibraryClosure
- catalogRefresher
- fcntl
- copyRegularAt
- dataRowRenderer
- runGuard
- btrfsLayoutFrom
- Ledger
- SanitizeText
- runMonitor
- process_vm_readv.c
- runForward
- runUninstall
- update.go
- CandidateDir
- GuardInodeKey
- os.FileMode
- render.go
- time.Time
- TempBinary
- daemon mode (fscrypt + whitelist lifecycle)
- FORWARD temporary -w/-b whitelist grants
- KernelDev
- configedit_test.go
- netinject.c
- NetGuard
- check-compatibility.sh
- BPF verifier 1M-insn complexity budget
- Inode-based (dev:ino) executable identity
- expandPlaceholders
- netModel
- IntegrationSuite
- Well-known bypass classes
- Vault
- jit_provenance.c
- inspector_probe.c
- ptrace_race.c
- NetGuard
- Build & sign release assets (reusable workflow)
- testLimiter
- tempGrantSetup
- net_tester/main.go
- NetworkGuardUseCase
- app-listener Terminal Demo Recording
- auth.go
- runOfflineEdit
- headerWidget
- IntegrationSuite
- IntegrationSuite
- Regenerate eBPF bindings in pinned builder image
- auditGroupCoverage
- Config
- WebSocket /ws client
- github.com/spf13/cobra.Command
- launch_probe.c
- install.sh
- mc_tag.h
- .step
- make test-integration (rootful Docker suite)
- OpenSystemPlaced
- supersede_probe.c
- bufio.Reader
- trace-app-libs.sh
- IntegrationSuite
- CheckBPFLSM
- Bun TMPDIR redirect with reserved .bun-* name
- Daemon self-protection (selfguards.go)
- github.com/Virgula0/app-listener
- network-monitor mode
- uninstall subcommand
- ResolveLibraryClosure
- buildTrustedSet
- os.File
- Guard
- Run
- vouchedDevs
- classifySupportError
- GlobReservations
- writeWithin
- daemon: permanent protection with encryption at rest
- IntegrationSuite
- inodeChain
- stripTempJournal
- newTestVetter
- is_event_type_allowed
- Installation
- secureResources
- graphify-refresh.sh
- TestMain
- .degradeWatchLoop
- find_test.go
- Candidate
- dialLiveSession
- IntegrationSuite
- mustSubkey
- bpfstats/main.go
- writeEvent
- Limitations
- Troubleshooting
- watchPattern
- fakeVault
- Daemon configuration
- edit-protected
- runLockdown
- launchKey
- Compatibility
- How it works
- parseConfig
- ParseEventsFlag
- bpfLsmListed
- findTopLevelWindow
- Binary trust
- guard: block file access by program
- network-guard: block network access by program
- launch_rule_at
- Development
- app-listener documentation
- Self-updating apps and the catalog
- network-monitor: see what a program talks to
- unlockUnderGuard
- CleanupStalePins
- .runEditorHarness
- monitor: see what touches a path
- kernel-prof.sh
- Vault
- ldConfParser
- copyXattrs
- golang.org/x/sys/unix.Timespec
- VerifyInstalledResourcesLocked
- MulticallError

## God Nodes (most connected - your core abstractions)
1. `Guard` - 96 edges
2. `Config` - 93 edges
3. `Load()` - 85 edges
4. `absPath()` - 78 edges
5. `IntegrationSuite` - 75 edges
6. `IntegrationSuite` - 66 edges
7. `writeConfig()` - 65 edges
8. `fileEditModel` - 60 edges
9. `User` - 56 edges
10. `GuardInodeKey` - 55 edges

## Surprising Connections (you probably didn't know these)
- `Well-known bypass classes` --references--> `ptrace_race POC`  [INFERRED]
  CLAUDE.md → integrationtests/exploits/README.md
- `lib_probe reserved library POC` --semantically_similar_to--> `Bun TMPDIR redirect with reserved .bun-* name`  [INFERRED] [semantically similar]
  integrationtests/exploits/README.md → README.md
- `swap_benign / swap_reader in-place binary swap POC` --semantically_similar_to--> `Binary ledger (/etc/app-listener/binaries.db)`  [INFERRED] [semantically similar]
  integrationtests/exploits/README.md → README.md
- `Well-known bypass classes` --references--> `btrfs_search POC`  [INFERRED]
  CLAUDE.md → integrationtests/exploits/README.md
- `Well-known bypass classes` --references--> `io_uring POC`  [INFERRED]
  CLAUDE.md → integrationtests/exploits/README.md

## Import Cycles
- None detected.

## Hyperedges (group relationships)
- **Shared bubbletea TUI layout across monitor and guard modes** — media_demo_monitor_tui, media_demo_guard_whitelist_tui, media_demo_guard_blacklist_tui, media_demo_tui_resource_footer [INFERRED 0.85]
- **Reserved-name trust defenses (trust object #3)** — integrationtests_exploits_readme_glob_plant, integrationtests_exploits_readme_lib_probe, readme_bun_tmpdir_redirect, claude_catalog_refresh [INFERRED 0.85]
- **Release build, sign and verify pipeline** — _github_workflows_release_prerelease_workflow, _github_workflows_promote_release_promote_release_workflow, _github_workflows__build_release_build_release_workflow, _github_workflows__build_release_ed25519_signing, readme_install_sh, readme_update [INFERRED 0.95]
- **Verifier budget measurement and gating** — _github_workflows_verifier_verifier_gate, claude_bpfstats, readme_daemon_check, claude_verifier_complexity_budget, _github_workflows__build_release_ebpf_regeneration [INFERRED 0.95]

## Communities (212 total, 42 thin omitted)

### Community 0 - "daemon/daemon.go"
Cohesion: 0.10
Nodes (16): Execute(), tempJournalRow, launchRule, trustEvent, btrfsIoctlFsInfoArgs, Check(), collectDiagnostics(), readSysctl() (+8 more)

### Community 1 - "testing.T"
Cohesion: 0.02
Nodes (126): plantBun(), TestOpenBunTmp_KeepsVettedDir(), discordRefreshFixture(), TestRefreshOnceAppliesAndReloads(), TestRefreshOnceKeepsOperatorLinesAndSkipsNoOps(), TestRefreshOnceRestoresTheFileWhenTheReloadFails(), TestRefreshOnceSkipsAConfigChangedSinceLoad(), TestDebouncerCoalescesAndBoundsWait() (+118 more)

### Community 2 - "IntegrationSuite"
Cohesion: 0.09
Nodes (10): configBinaryLines(), IntegrationSuite, parseDaemonEvents(), probeExeAs(), procGateMatrix(), procReadProbe(), rawExec(), shQuote() (+2 more)

### Community 3 - "filevault.go"
Cohesion: 0.23
Nodes (16): classifyRegularFileTarget(), clearRecovery(), deriveFileVaultSubkey(), EnsureRecoverySidecarPlaceholder(), fileVaultRecoverPath(), Vault, openFileVaultWithMasterKey(), recoverFileInPlace() (+8 more)

### Community 4 - "daemonconfig_test.go"
Cohesion: 0.08
Nodes (72): Load(), ResolvePendingPaths(), hasBinary(), hasPendingBinary(), TestElectronAppsBlockParsed(), TestElectronAppsBlockRejectsMalformedLines(), TestEncryptionGroupsSkipsLibDirs(), TestInspectorsBlockBounded() (+64 more)

### Community 5 - "networkguard.bpf.c"
Cohesion: 0.10
Nodes (38): admit_refusal(), check_watched(), code_vouched(), emit_gate(), env_loader(), env_step(), exe_inode_of(), fill_current_image() (+30 more)

### Community 6 - "monitor.bpf.c"
Cohesion: 0.15
Nodes (31): emit_event(), emit_event_kern(), fill_path(), read_path(), trace_do_faccessat(), trace_do_sendfile(), trace_do_splice(), trace_do_splice_direct() (+23 more)

### Community 7 - "User"
Cohesion: 0.11
Nodes (40): trustManager, rcTarget, BunTmpDir(), catalogEntryMatch(), findCatalogRoot(), findCatalogWatchSubPath(), isInsidePath(), ResolveCatalogEntry() (+32 more)

### Community 8 - "guardUnitTest"
Cohesion: 0.05
Nodes (26): addErrKind, guardUnitTest, ReplacementCheck, BinaryEntry, BinariesSummary(), classifyAddErr(), addRecorder(), copySelf() (+18 more)

### Community 9 - "mkdirs"
Cohesion: 0.14
Nodes (35): newWatchPattern(), newFakeInotify(), TestWatchSetBinaryWriteTriggersNotMkdir(), TestWatchSetFloodIsCapped(), TestWatchSetLibDirAppearanceTriggers(), TestWatchSetMissingRootWatchesDeepestAncestor(), TestWatchSetOverflowRescansAndRebuilds(), TestWatchSetSymlinkCreateTriggers() (+27 more)

### Community 10 - "Hash"
Cohesion: 0.18
Nodes (20): Hash(), HashFileExists(), LoadHashFile(), RemoveHashFile(), hashFileInode(), TestHashFileLifecycle(), TestRemoveHashFileInPlace(), TestValidatePassword() (+12 more)

### Community 11 - "planUpdaters"
Cohesion: 0.16
Nodes (16): hardLinked(), libBinaries(), mayUpdate(), planUpdaters(), resourceBinaries(), resourceLinks(), rules(), TestGeneralToolsToWarn_HardLinkIsNoUpdaterButNoWarning() (+8 more)

### Community 12 - "codeedit.go"
Cohesion: 0.11
Nodes (41): nextConfig(), class, classOf(), dedent(), Edit(), leadingIndent(), lines(), moveTo() (+33 more)

### Community 13 - "IntegrationSuite"
Cohesion: 0.16
Nodes (7): netMonitorEvent, IntegrationSuite, fmtEvents(), IntegrationSuite, netMonitorDeltaEvents(), netMonitorHasTypes(), parseNetMonitorEvents()

### Community 14 - "RunCmd"
Cohesion: 0.07
Nodes (36): reloadDaemonForPasswordChange(), applyLiveRefresh(), buildBinaryIfNeeded(), installBinaryOnly(), preflightDeployedBinary(), runUpdateCatalogOnly(), softenAutomatedRefreshErr(), TestApplyLiveRefreshDeliversReload() (+28 more)

### Community 15 - "Vault"
Cohesion: 0.10
Nodes (23): lockOneRoot(), lockRootRecovering(), relockStaleVaults(), deprovisionKind, isFileVaultCiphertext(), isRegularFileTarget(), classifyDeprovision(), classifyDeprovisionErr() (+15 more)

### Community 16 - "IntegrationSuite"
Cohesion: 0.11
Nodes (9): attrOpCase, eventsForPath(), IntegrationSuite, guardDeltaEvents(), guardEventTypes(), le64HexKey(), parseGuardEvents(), guardEvent (+1 more)

### Community 17 - "guard_trust.bpf.c"
Cohesion: 0.11
Nodes (51): caller_updates(), current_exe_flags(), current_exe_inode(), current_exe_key(), current_writer_bits(), dentry_is_dir(), dentry_parent(), deny_if_protected() (+43 more)

### Community 18 - "RefreshSection"
Cohesion: 0.12
Nodes (18): RefreshOptions, isLibDirectiveText(), ParseSectionHeaderPath(), unquotePath(), IsLibraryDirective(), admitNew(), keepExisting(), LiveEmptyWhitelistRejected() (+10 more)

### Community 19 - "fileEditModel"
Cohesion: 0.10
Nodes (9): clipLabel(), humanSize(), clamp(), sortNodes(), chownUser, fileEditMode, fileEditModel, fileKind (+1 more)

### Community 20 - "networkmonitor.bpf.c"
Cohesion: 0.31
Nodes (16): BPF_KRETPROBE(), emit_event(), get_netns(), get_socket_proto(), is_watched_binary(), read_inet_addr(), trace_accept(), trace_accept4() (+8 more)

### Community 21 - "usecase/daemon_test.go"
Cohesion: 0.20
Nodes (43): NewDaemonUseCase(), newFakeVault(), resource(), startDaemon(), TestDaemonUseCaseConstructorMismatch(), TestDaemonUseCaseEventsMerged(), TestDaemonUseCaseGrantEditAccess(), TestDaemonUseCaseGroupedEncryptionRootsDeduplicated() (+35 more)

### Community 22 - "Resource"
Cohesion: 0.08
Nodes (12): ownerPaths(), Resource, GuardRepository, concurrencyLimit(), encryptionRootOf(), DaemonUseCase, partitionEncryptionRoots(), TestPartitionEncryptionRootsSplitsDirsAndFiles() (+4 more)

### Community 23 - "RawTarget"
Cohesion: 0.27
Nodes (10): BuildEBPFTargets(), RawTarget, isSubDir(), pathWithinDir(), ResolveTargets(), validatePair(), ValidateTargets(), WarnIgnoredFlags() (+2 more)

### Community 24 - "runNetworkGuard"
Cohesion: 0.15
Nodes (18): ParseNetEventsFlag(), computeGuardBinaries(), confirmUnsafe(), discoverInfra(), resolveGuardMode(), runNetworkGuard(), runTUI(), DiscoverInfraBinaries() (+10 more)

### Community 25 - "guard.bpf.c"
Cohesion: 0.12
Nodes (46): add_inode_to_guard(), chmod_only_drops_write(), count_degrade(), discover_guarded_parent(), evict_inode_from_guard(), fill_path(), get_dentry_from_path(), get_inode_from_path() (+38 more)

### Community 26 - "NetworkMonitor"
Cohesion: 0.10
Nodes (14): NetBpfEvent, NetEvent, NetEventType, Cstr(), FormatAddr(), NetEventTypes(), Ntohs(), ParseNetEventType() (+6 more)

### Community 27 - "monitorUnitTest"
Cohesion: 0.10
Nodes (3): newPathCache(), dirTarget(), monitorUnitTest

### Community 28 - "runDaemon"
Cohesion: 0.16
Nodes (10): catchLifecycleSignals(), makeReloadHandler(), notifySystemdReady(), runDaemon(), startCatalogRefresh(), startPprof(), writePidFile(), mergeDaemonEvents() (+2 more)

### Community 29 - "buildOneGuard"
Cohesion: 0.09
Nodes (37): backingDeviceUnion(), buildConcurrency(), buildGuards(), buildOneGuard(), eventFilterOptions(), reloadOnce(), reloadSlotsNeeded(), startGuardedDaemon() (+29 more)

### Community 30 - "Monitor"
Cohesion: 0.11
Nodes (5): evalSymlinksOrEmpty(), NewMonitor(), Monitor, pathCache, pathCacheEntry

### Community 31 - "serve.go"
Cohesion: 0.07
Nodes (35): basicAuth(), configureWebSocket(), decodeEndsAtEOF(), isInteractiveTerminal(), NewEventFanout(), newServeListener(), newServeServer(), parseResize() (+27 more)

### Community 32 - "engine"
Cohesion: 0.07
Nodes (18): inodeBatchFn, TestNearestRootClaims(), TestTieRank(), appendInodeRow(), batchInodeRows(), engine, inodeRowsOf(), nearestRootClaims() (+10 more)

### Community 34 - "fileedit_model_test.go"
Cohesion: 0.10
Nodes (38): TestRuntimeClassMarkers(), TestFilterExistingWhitelistSymlinkEscape(), AtomicWriteAt(), OpenRegularNoFollow(), defaultNewFileMode(), fileModeToUnixMode(), isBinaryContent(), isOctalMode() (+30 more)

### Community 35 - "daemonModel"
Cohesion: 0.27
Nodes (6): listenForDaemonEvents(), tickDaemonStats(), daemonEventLine, daemonEventMsg, daemonModel, DaemonResourceInfo

### Community 36 - "watchGroup"
Cohesion: 0.11
Nodes (35): watchGroup, addResource(), applyAllowLib(), applyConfigLine(), applyDirective(), applyElectronApp(), applyInspector(), applyLibBinary() (+27 more)

### Community 37 - "filevault_test.go"
Cohesion: 0.25
Nodes (19): looksLikeFileVaultRecord(), inode(), TestIsEncryptedDispatchesToFileVaultForRegularFiles(), TestIsProvisionedForFile(), TestLockUnlockFileInPlaceIdempotent(), TestUnlockFileInPlaceRecoversFromInterruptedTransform(), TestUnlockFileInPlaceWrongKeyFails(), TestUnlockLockFileInPlacePreservesInode() (+11 more)

### Community 38 - "GuardEvent"
Cohesion: 0.12
Nodes (12): bpfGuardEvent, GuardEvent, commMatchesGuardedBinary(), decodeGuardEvent(), fsGateLabel(), logBacklogDenial(), parseGuardEvent(), processGateLabel() (+4 more)

### Community 39 - "guardSpec"
Cohesion: 0.25
Nodes (10): guardSpec(), TestGuardSpecKeepsIoctlCompatOnlyWhereKernelHasIt(), verifyCollection(), VerifyLoad(), InertProgram(), KernelHasFunc(), TestKernelHasFunc(), trustSpec() (+2 more)

### Community 40 - "github.com/charmbracelet/bubbles/textarea.Model"
Cohesion: 0.16
Nodes (10): findButton, span, column(), findAll(), Finder, lowerRunes(), newFindInput(), overlay() (+2 more)

### Community 41 - "runNetworkMonitor"
Cohesion: 0.14
Nodes (9): runHeadless(), runHeadless(), runNetworkMonitor(), runTUI(), ProtocolString(), NetworkMonitorRepository, NewNetModel(), NetworkMonitorUseCase (+1 more)

### Community 42 - "absPath"
Cohesion: 0.08
Nodes (4): IntegrationSuite, inNosuidNamespace(), IntegrationSuite, absPath()

### Community 44 - "IntegrationSuite"
Cohesion: 0.17
Nodes (8): guardBinaryFlag(), infraContainerPaths(), IntegrationSuite, IntegrationSuite, guardNetTypesForComm(), netGuardBlockedEventCount(), netGuardHasBlockedEvent(), netGuardTail()

### Community 45 - "SelectOrphans"
Cohesion: 0.24
Nodes (16): appProtectorSet(), CleanOrphans(), isAppProtector(), listPolicies(), listProtectors(), orphanedPolicies(), orphanedProtectors(), SelectOrphans() (+8 more)

### Community 48 - ".syncSetsLocked"
Cohesion: 0.13
Nodes (16): resInfo, taintOwner, nextLabel(), TestDeliverLabelsCrossResourceProcessGates(), engine, isSubset(), planMemberRows(), planTaintOwners() (+8 more)

### Community 49 - "TrustGuard"
Cohesion: 0.10
Nodes (10): statFunc, trustHook, trustRows, setSupersedeTrust(), cStr(), resolveBits(), statKey(), TrustGuard (+2 more)

### Community 50 - "tempRow"
Cohesion: 0.17
Nodes (7): exeRow, tempGrant, tempMaskOp, tempRow, Guard, liveTrust(), planTempBlock()

### Community 51 - "__always_inline"
Cohesion: 0.11
Nodes (27): exe_is_superseded(), exe_refused(), ctx_ptr(), inode_dev(), sb_dev(), btrfs_copy_gate(), check_and_emit_args(), code_suspect() (+19 more)

### Community 52 - "IntegrationSuite"
Cohesion: 0.06
Nodes (15): IntegrationSuite, decodeMulticall(), IntegrationSuite, multicallBypasses(), multicallConfig(), uutilsStrays(), vettedGroupIndex(), exploitTest (+7 more)

### Community 53 - "MountVouch"
Cohesion: 0.08
Nodes (18): checkBtrfsKeys(), ensureBtrfsLayout(), deleteResKeys(), syncU32Map(), syncMap(), deleteKeys(), MountVouch, TrustGuard (+10 more)

### Community 54 - "github.com/charmbracelet/bubbletea.Cmd"
Cohesion: 0.07
Nodes (15): diffModel, editorModel, listenForGuardEvents(), tickGuardStats(), protoString(), renderBullets(), TestRenderBulletsVerbatimAndIndented(), guardEventMsg (+7 more)

### Community 56 - "highlight.go"
Cohesion: 0.23
Nodes (5): fg(), Highlighter, lexerFor(), shebangLexer(), tokenClassOf()

### Community 57 - "sync.Mutex"
Cohesion: 0.20
Nodes (5): controlManager, controlManager, newControlManager(), startControlManager(), configEditor

### Community 58 - "GenerateConf"
Cohesion: 0.07
Nodes (51): appendSectionsAndEdit(), TestDiffMergePreservesExistingAndParses(), editConfig(), libraryBlocksFromCandidates(), runConfigEditor(), sectionsFromCandidates(), discordCatalogEntry(), groupedDiscordConf() (+43 more)

### Community 59 - "EventType"
Cohesion: 0.09
Nodes (14): BpfEvent, EventType, eventUnitTest, FileEvent, parseEvents(), TestAddBinaryActionsRefusesBeforeWriting(), TestCheckBinaryEventsReadOnlyRejectsRestriction(), Guard (+6 more)

### Community 60 - "fileedit_symlink_test.go"
Cohesion: 0.23
Nodes (11): mustRead(), mustWrite(), swappedParent(), TestChmodRefusesSymlinkedParent(), TestCreateEntryRefusesPlantedSymlink(), TestCreateEntryRefusesSymlinkedParent(), TestDeleteRefusesSymlinkedParent(), TestSaveFollowsPinnedVaultNotSwappedAncestor() (+3 more)

### Community 61 - "ListUsers"
Cohesion: 0.10
Nodes (28): askBunTmpdirUsers(), bunEntryConfigured(), bunLaunchersFor(), gatherPerUserSetup(), setupBunTmpdir(), applyDiffAdditions(), configPaths(), restoreDaemonAfterDiffAbort() (+20 more)

### Community 62 - "exe_supersede.h"
Cohesion: 0.19
Nodes (17): exe_last_link_key(), exe_note_freed(), exe_note_supersede(), exe_seq_now(), exe_stamp_exec(), exe_stamp_fork(), exe_stamp_free(), exe_stamp_lost() (+9 more)

### Community 64 - "errno"
Cohesion: 0.21
Nodes (5): put(), wait_then_read(), dump_via_debugfs(), main(), probe_open()

### Community 65 - "pinstate.go"
Cohesion: 0.29
Nodes (14): pinOwnerLikelyAlive(), readPinState(), recoverPinState(), pinStateInode(), TestPinOwnerLikelyAlive(), TestPinOwnerLikelyAliveDeadPID(), TestPinStateRoundTrip(), TestReadPinStateMalformed() (+6 more)

### Community 66 - "controlServer"
Cohesion: 0.21
Nodes (6): controlServer, readClientLines(), startControlServer(), streamEvents(), editControlSession, grantRequest

### Community 67 - "time.Duration"
Cohesion: 0.21
Nodes (4): newTestSession(), TestRunSessionIdleTimeoutResetsOnPing(), TestRunSessionStreamsEventsAsActivity(), debouncer

### Community 68 - "LibraryClosure"
Cohesion: 0.14
Nodes (15): defaultLibDirs(), elfInterp(), isScript(), ldSoConfDirs(), LibraryClosure(), resolveSoname(), rootOwnedInode(), rootOwnedSafe() (+7 more)

### Community 69 - "catalogRefresher"
Cohesion: 0.13
Nodes (11): drain(), newCatalogRefresher(), readEvents(), sameConfig(), swapFile(), awaitStartupOrSignal(), TestAwaitStartupOrSignalAbortsAndLocksDown(), TestAwaitStartupOrSignalAbortWithFailedStartup() (+3 more)

### Community 70 - "fcntl"
Cohesion: 0.13
Nodes (4): denied(), main(), denied(), main()

### Community 71 - "copyRegularAt"
Cohesion: 0.31
Nodes (10): copyDirContents(), copyEntry(), copyRegularAt(), copySymlinkAt(), copyTreeRoot(), fchownFromInfo(), futimesFromInfo(), readDirNames() (+2 more)

### Community 73 - "runGuard"
Cohesion: 0.13
Nodes (17): displayEntries(), holdForward(), holdHeadless(), pingInterval(), chooseResource(), runLiveEdit(), HeadlessLine(), resolveGuardConfig() (+9 more)

### Community 74 - "btrfsLayoutFrom"
Cohesion: 0.15
Nodes (11): BtrfsLayout, fakeTypes, typeSource, btrfsLayoutFrom(), findMember(), memberOffset(), ResolveBtrfsLayout(), btrfsTypes() (+3 more)

### Community 75 - "Ledger"
Cohesion: 0.08
Nodes (23): logRefreshChange(), TestLogRefreshChange(), grantNewBinaries(), candidates(), confirmable(), confirmAll(), describe(), hashNow() (+15 more)

### Community 76 - "SanitizeText"
Cohesion: 0.19
Nodes (7): multicallAdmissible(), newBinaryVetter(), openBinaryVetter(), drainEphemeral(), binaryVetter, vetDecision, SanitizeText()

### Community 78 - "runMonitor"
Cohesion: 0.11
Nodes (14): CheckEBPF(), MakeDisplayPaths(), newPinGeneration(), prepareDaemonStart(), runHeadless(), runMonitor(), runTUI(), PrintLogo() (+6 more)

### Community 79 - "process_vm_readv.c"
Cohesion: 0.40
Nodes (6): dump_memory(), find_all_rw_regions(), find_heap(), find_victim_by_fd(), main(), spawn_less_pty()

### Community 80 - "runForward"
Cohesion: 0.40
Nodes (6): chooseResources(), confirmUserPlaced(), eventsLabel(), forwardBinaries(), normalizedEvents(), runForward()

### Community 81 - "runUninstall"
Cohesion: 0.09
Nodes (34): parseMaintenanceFlags(), runMaintenanceMode(), TestValidateMaintenanceFlags(), validateMaintenanceFlags(), deletePostBackups(), restoreBackups(), offerBackupCleanup(), runUninstall() (+26 more)

### Community 82 - "update.go"
Cohesion: 0.05
Nodes (53): applyUpdate(), assetsFor(), checkUpdatePreconditions(), compareStableVersions(), downloadAndVerify(), downloadFile(), downloadReleaseFiles(), fetchReleases() (+45 more)

### Community 83 - "CandidateDir"
Cohesion: 0.10
Nodes (20): refreshAdmits(), TestRefreshAdmitsHomeMatchesAndRootPlacedNamesOnly(), inHome(), entryWriters(), reserveChain(), underResource(), globBuilder, ParseGlobName() (+12 more)

### Community 84 - "GuardInodeKey"
Cohesion: 0.06
Nodes (27): refusedReplacements, rootHandle, engine, Inspector, SetInspectors(), engine, supersedeUnnamed(), appletTag() (+19 more)

### Community 85 - "os.FileMode"
Cohesion: 0.10
Nodes (29): TestDaemonUnitWithMetadataOutput(), updateCatalogConfig(), daemonUnitWithMetadataOutput(), dupFD(), installFile(), installFileAs(), installServices(), installSSHAgent() (+21 more)

### Community 86 - "render.go"
Cohesion: 0.15
Nodes (13): lineStyles, viewportState, baseStyle(), Highlighter, gutterLabel(), lineLen(), markStyle(), repeatSpaces() (+5 more)

### Community 87 - "time.Time"
Cohesion: 0.23
Nodes (7): isSymlink(), newWatchSet(), fakeInotify, watchDir, watchPos, watchSet, dirEntry

### Community 88 - "TempBinary"
Cohesion: 0.07
Nodes (23): fakeDaemon, BinaryStat, binaryStatOf(), TempBinary, TemporaryGrant, TempRule, loadPinnedExeMaps(), ResolveTempBinary() (+15 more)

### Community 89 - "daemon mode (fscrypt + whitelist lifecycle)"
Cohesion: 0.18
Nodes (10): Daemon catalog refresh (catalogrefresh.go / catalogwatch.go), exe_supersede.h superseded-key tracking, glob_plant reserved glob name POC, launch_probe risky-launch POC, preload_wait code-suspect POC, supersede_probe reused-inode POC, Daemon boot ordering (Type=notify, Before=user sessions), daemon mode (fscrypt + whitelist lifecycle) (+2 more)

### Community 90 - "FORWARD temporary -w/-b whitelist grants"
Cohesion: 0.21
Nodes (9): CONFIG post-AUTH daemon.conf editing (configedit.go), FORWARD temporary -w/-b whitelist grants, Daemon self-key (GUARD_ALLOW_ROOT) bypass, swap_benign / swap_reader in-place binary swap POC, edit-protected --edit-config, edit-protected subcommand, edit-protected live mode, edit-protected --forward (+1 more)

### Community 91 - "KernelDev"
Cohesion: 0.10
Nodes (26): fdStat, physicalParent(), statFD(), KernelDev(), LstatInode(), BtrfsDevices(), btrfsSysfsDevices(), TestBtrfsSysfsDevices() (+18 more)

### Community 92 - "configedit_test.go"
Cohesion: 0.47
Nodes (9): newTestConfigEditor(), readFile(), TestConfigEditorApply(), TestConfigEditorRefusals(), TestConfigEditorRefusesNewMulticallLine(), TestConfigEditorRestoresOnFailedReload(), TestRunConfigEditOverTheSocket(), TestRunConfigEditRefusedDuringALiveSession() (+1 more)

### Community 93 - "netinject.c"
Cohesion: 0.17
Nodes (10): attach_child(), await_exec(), do_connect(), errname(), main(), sleep_ms(), traceme(), do_connect() (+2 more)

### Community 95 - "NetGuard"
Cohesion: 0.13
Nodes (7): foreignPinnedLink(), eventTypeKey(), NetGuard, NetGuardEvent, mcNameKey(), statInodeKey(), netGuardEventLine

### Community 96 - "check-compatibility.sh"
Cohesion: 0.38
Nodes (11): fail(), load_config(), lsm_bpf_instructions(), note(), pass(), read_config(), require_config(), section() (+3 more)

### Community 97 - "BPF verifier 1M-insn complexity budget"
Cohesion: 0.20
Nodes (8): PR vs base verifier cost comparison, Verifier gate workflow, bpf_loop() helper, tools/bpfstats verifier cost tool, static __noinline BPF-to-BPF helpers with args struct, rename-over-guarded-file bypass, xres_move cross-resource move POC, Kernel floor (5.8 / 5.10 / 5.17)

### Community 98 - "Inode-based (dev:ino) executable identity"
Cohesion: 0.18
Nodes (7): CheckBPFLSM preflight, internal/infrastructure shared eBPF kernel, Ports & adapters layering (cmd -> usecase -> repository -> engines), jit_provenance ledger (guard_jit_origin), daemon --check / --verifier-only preflight, guard mode (LSM whitelist/blacklist), network-guard mode

### Community 99 - "expandPlaceholders"
Cohesion: 0.18
Nodes (8): symlinkTargetGlob(), expandPlaceholders(), CandidateDir, CandidateDir, LibDirGlob, TrustGlob, reservedGlobs(), splitTrustGlob()

### Community 100 - "netModel"
Cohesion: 0.06
Nodes (33): Stats, parseKBValue(), parseStat(), parseStatus(), Read(), TestReadConcurrentRace(), TestStatsFieldsNonZero(), TestStatsValuesMonotonic() (+25 more)

### Community 102 - "Well-known bypass classes"
Cohesion: 0.24
Nodes (9): btrfs copy-ioctl gate (file_ioctl / file_ioctl_compat), Well-known bypass classes, btrfs_search POC, io_uring POC, open_by_handle_at POC, process_vm_readv POC (remaining monitor gap), raw_block_device POC, statonly stat/statx metadata leak (+1 more)

### Community 103 - "Vault"
Cohesion: 0.19
Nodes (9): migrationTarget, writeTempWithMetadata(), classifyMigrationTarget(), Vault, readClassifiedFile(), stampCopyMetadata(), CopyTreeWithProgress(), TestCopyTreeProgress() (+1 more)

### Community 104 - "jit_provenance.c"
Cohesion: 0.47
Nodes (8): copy_into(), main(), make_temp_copy(), mode_hold(), mode_memfd(), mode_passfd(), mode_self(), report_dlopen()

### Community 105 - "inspector_probe.c"
Cohesion: 0.56
Nodes (8): check_attach(), check_exe(), check_mem(), check_root(), check_vmread(), main(), report(), traced_exec()

### Community 106 - "ptrace_race.c"
Cohesion: 0.31
Nodes (5): main(), peek_buffer(), run_tracer(), run_victim(), wait_for_file()

### Community 108 - "Build & sign release assets (reusable workflow)"
Cohesion: 0.31
Nodes (8): Build & sign release assets (reusable workflow), app-listener-gui best-effort build, Promote to stable release workflow, Release workflow (pre-release per push to main), Release promote flow, make check-compatibility, scripts/install.sh one-line installer, update self-update subcommand

### Community 109 - "testLimiter"
Cohesion: 0.26
Nodes (11): newGateLogLimiter(), gateEvent(), TestGateLogLimiterFlushesOnShutdown(), TestGateLogLimiterFoldsAcrossThreadsAndPids(), TestGateLogLimiterFoldsReadRepeats(), TestGateLogLimiterKeysOnTarget(), TestGateLogLimiterNeverFoldsWhatMatters(), testLimiter() (+3 more)

### Community 110 - "tempGrantSetup"
Cohesion: 0.21
Nodes (10): tempJournalRows(), journalTo(), tempBinary(), tempGrantSetup(), TestGrantTemporaryAccessJournalsBeforeApplyAndClearsAfterRevoke(), TestGrantTemporaryAccessRefusals(), TestGrantTemporaryAccessRevokesOnInPlaceRewrite(), TestGrantTemporaryAccessRollsBackOnApplyFailure() (+2 more)

### Community 111 - "net_tester/main.go"
Cohesion: 0.21
Nodes (7): connFD(), tcpClient(), tcpServer(), tcpServerDelayed(), udpRecvLoop(), udpSendLoop(), waitForMarker()

### Community 112 - "NetworkGuardUseCase"
Cohesion: 0.15
Nodes (11): NetworkGuardRepository, NetworkGuardUseCase, NewNetworkGuardUseCase(), newFakeNetworkGuardRepo(), newFakeNetworkMonitorRepo(), TestGuardUseCaseLifecycle(), TestMonitorUseCaseLifecycle(), TestNetworkGuardUseCaseLifecycle() (+3 more)

### Community 113 - "app-listener Terminal Demo Recording"
Cohesion: 0.39
Nodes (8): BPF Build Pipeline (bpf2go + go build), CLI Subcommand Help (guard, install, monitor, network-guard, network-monitor), app-listener Terminal Demo Recording, Guard TUI Blacklist Mode (eBPF LSM), Guard TUI Whitelist Mode (eBPF LSM), File System Monitor TUI (eBPF), Startup Banner and eBPF Availability Check, TUI Resource Footer (Mem / CPU / Evt/s)

### Community 114 - "auth.go"
Cohesion: 0.20
Nodes (7): OriginOf(), parseHash(), TestHashVerifyRoundTrip(), TestVerifyMalformed(), Verify(), Origin, parsedHash

### Community 115 - "runOfflineEdit"
Cohesion: 0.15
Nodes (13): TestPickEncryptedDirSingleSkipsPicker(), pickEncryptedDir(), runOfflineEdit(), askDecrypt(), decryptDirectories(), decryptStep(), pickDirsToDecrypt(), TestDecryptDirectoriesEmpty() (+5 more)

### Community 116 - "headerWidget"
Cohesion: 0.18
Nodes (6): columnState, dataRow, headerWidget, newColumnState(), newDataRow(), newHeaderWidget()

### Community 119 - "Regenerate eBPF bindings in pinned builder image"
Cohesion: 0.33
Nodes (5): CI workflow (lint + test on non-draft PRs), CI lint job (make lint), CI test job (make test), golangci-lint v2 configuration, Issue #45 (clang 14 guard_path_rename over budget)

### Community 120 - "auditGroupCoverage"
Cohesion: 0.18
Nodes (13): auditAfterEdit(), auditAfterEditWithConfig(), auditEntry(), auditGroupCoverage(), auditTree(), coveredByWatch(), findResource(), groupWatchPaths() (+5 more)

### Community 121 - "Config"
Cohesion: 0.17
Nodes (19): ServeConfig, newDaemonModel(), runDaemonUI(), runHeadless(), runServedTUI(), runTUI(), admitEvent(), foldable() (+11 more)

### Community 122 - "WebSocket /ws client"
Cohesion: 0.29
Nodes (3): AppListener Shared Session page, WebSocket /ws client, --serve WebSocket mirroring

### Community 123 - "github.com/spf13/cobra.Command"
Cohesion: 0.08
Nodes (33): AddServeFlags(), credentialFlagState(), isLoopbackHost(), ParseServeFlags(), requestedBoolFlag(), validateEnabledServe(), validateServeAddress(), TestParseServeFlagsRequiresCredentials() (+25 more)

### Community 124 - "launch_probe.c"
Cohesion: 0.53
Nodes (4): copy(), is_verb(), main(), usage()

### Community 125 - "install.sh"
Cohesion: 0.57
Nodes (6): die(), download(), info(), install.sh script, usage(), warn()

### Community 126 - "mc_tag.h"
Cohesion: 0.24
Nodes (10): mc_attest(), mc_basename(), mc_leader_start(), mc_stamp_fork(), mc_stamp_free(), mc_tag_bits(), netg_exec_applet(), netg_task_free() (+2 more)

### Community 127 - ".step"
Cohesion: 0.20
Nodes (11): placedWalk, confineBelow(), inUserHome(), GlobSystemPlaced(), placedDir(), placedMatches(), readlinkFd(), rootOwnedStat() (+3 more)

### Community 128 - "make test-integration (rootful Docker suite)"
Cohesion: 0.40
Nodes (4): make test-integration (rootful Docker suite), app-listener docker-compose service, ptrace_race POC, Docker usage (docker compose)

### Community 129 - "OpenSystemPlaced"
Cohesion: 0.14
Nodes (14): loadDaemonConfig(), resolveConfigPath(), setConfinedHomes(), TestLoadDaemonConfigEmptyResourcesIsCriticalStartup(), TestLoadDaemonConfigMissingFileIsCriticalStartup(), fileKey(), TestSystemFileRefusesAnotherInode(), TestSystemFileRefusesUserPlacedName() (+6 more)

### Community 130 - "supersede_probe.c"
Cohesion: 0.90
Nodes (4): answer(), main(), read_secret(), run_child()

### Community 131 - "bufio.Reader"
Cohesion: 0.18
Nodes (12): controlServer, readConfigPut(), TestReadGrantRequestConfig(), readAuthRequest(), readForwardRequest(), readGrantRequest(), readPathLines(), TestReadAuthRequest() (+4 more)

### Community 132 - "trace-app-libs.sh"
Cohesion: 0.60
Nodes (3): conf_directive_values(), list_conf_binaries(), trace-app-libs.sh script

### Community 134 - "CheckBPFLSM"
Cohesion: 0.17
Nodes (13): CheckBPFLSM(), confirmMasterKeyOverwrite(), runBPFCheck(), runGenKey(), runOneShotMode(), ensureInstalledBinary(), ensureMasterKey(), prepareInstallation() (+5 more)

### Community 135 - "Bun TMPDIR redirect with reserved .bun-* name"
Cohesion: 0.50
Nodes (4): lib_probe reserved library POC, Bun TMPDIR redirect with reserved .bun-* name, install --diff-catalog, install TUI wizard

### Community 142 - "ResolveLibraryClosure"
Cohesion: 0.15
Nodes (13): fileExists(), hostView(), LdSoPreloadPaths(), mountVouchesOwnership(), ResolveLibraryClosure(), skipWithoutHostRoot(), TestHostView(), TestMountVouchesOwnershipJudgesHostNamespace() (+5 more)

### Community 143 - "buildTrustedSet"
Cohesion: 0.21
Nodes (13): inspectorPaths(), TestInspectorPathsOnlyRootPlaced(), buildTrustedSet(), TestBuildTrustedSetExcludesLibraryClosure(), closureRejections(), inGuardedTree(), keys(), TestWarnUntrustedLibs_OnlyWhatTheKernelRefuses() (+5 more)

### Community 144 - "os.File"
Cohesion: 0.06
Nodes (46): checkFreshDir(), fileID(), openBunTmp(), openDirNoFollow(), renewDir(), TestOpenBunTmp_ReplacesUnvettedDir(), TestOpenBunTmp_SymlinkReplacedTargetUntouched(), bunTmp (+38 more)

### Community 146 - "Guard"
Cohesion: 0.05
Nodes (25): binaryVerifyState, deferredBinary, Guard, Mode, sharedHashEntry, VettedInode, openRoot(), TestWithChmodDropWriteNeedsReadOnly() (+17 more)

### Community 148 - "vouchedDevs"
Cohesion: 0.33
Nodes (6): isFuseType(), parseMajorMinor(), TestParseMajorMinorMatchesStatEncoding(), TestVouchedDevs(), TestVouchedDevsFailsClosed(), vouchedDevs()

### Community 149 - "classifySupportError"
Cohesion: 0.18
Nodes (10): classifySetupError(), classifySupportError(), TestClassifySetupErrorGeneric(), TestClassifySetupErrorNotSetup(), TestClassifySetupErrorNotSupported(), TestClassifySupportErrorEncryptionNotEnabled(), TestClassifySupportErrorEncryptionNotEnabledF2fs(), TestClassifySupportErrorGeneric() (+2 more)

### Community 150 - "GlobReservations"
Cohesion: 0.21
Nodes (7): GlobChild, compileGlobNames(), globText(), GlobReservations, TrustGuard, resolveChildren(), TestCompileGlobNames_RepeatedPatternOrsBits()

### Community 151 - "writeWithin"
Cohesion: 0.22
Nodes (9): atomicWriteAt(), descendNoFollow(), refuseExistingSymlink(), TestWriteWithin_PlainFile(), TestWriteWithin_RefusesSymlinkDirComponent(), TestWriteWithin_RefusesSymlinkTarget(), TestWriteWithin_RejectsEscape(), writeAndSync() (+1 more)

### Community 152 - "daemon: permanent protection with encryption at rest"
Cohesion: 0.17
Nodes (12): Boot ordering, daemon: permanent protection with encryption at rest, Filesystem-wide gates, Flags, Lifecycle, Process inspectors and taint, Quiet metadata denials, Related (+4 more)

### Community 154 - "inodeChain"
Cohesion: 0.29
Nodes (12): inodeChain(), isProcFD(), mountinfoShowsSbRoot(), mountShowsSbRoot(), keyOf(), mkTree(), TestInodeChainDirIncludesItself(), TestInodeChainFollow() (+4 more)

### Community 155 - "stripTempJournal"
Cohesion: 0.22
Nodes (11): ensurePinStateFilePlaceholder(), ensurePlaceholder(), writeInPlace(), bootstrapSelfProtectPlaceholders(), ensureHashFilePlaceholder(), readTempJournal(), stripTempJournal(), TestTempJournalRoundTripInPlace() (+3 more)

### Community 156 - "newTestVetter"
Cohesion: 0.42
Nodes (9): newTestVetter(), testBinary(), TestVetterBootstrapRecordsThenChecks(), TestVetterLinkKeyedOnLink(), TestVetterRecordLiveUnderLinkLines(), TestVetterRefusesNewLineAfterBootstrap(), TestVetterResolverOnlyApproved(), TestVetterSwapWithinGenerationLosesRights() (+1 more)

### Community 157 - "is_event_type_allowed"
Cohesion: 0.44
Nodes (10): emit_event(), guard_net_socket_bind(), guard_net_socket_connect(), guard_net_socket_listen(), guard_net_socket_recvmsg(), guard_net_socket_sendmsg(), is_event_type_allowed(), is_socket_guarded() (+2 more)

### Community 158 - "Installation"
Cohesion: 0.18
Nodes (11): Build from source, Bun-based apps (opencode), Docker, fscrypt prerequisites, Install flags, Installation, One-line installer, SSH agent (+3 more)

### Community 159 - "secureResources"
Cohesion: 0.15
Nodes (13): encryptDirectories(), secureResources(), verifyEncryptionState(), collectFilesystemPrereqs(), confirmRunPrereq(), resolveFilesystemPrereqs(), askEncryption(), TestAskEncryptionSkipsNeedEncryptionFalse() (+5 more)

### Community 161 - "TestMain"
Cohesion: 0.50
Nodes (3): TestMain(), helperChild(), TestMain()

### Community 163 - "find_test.go"
Cohesion: 0.18
Nodes (20): changelogText(), newChangelogModel(), showChangelog(), TestChangelogModelOtherKeysDoNotQuit(), TestChangelogModelQuitKeys(), TestChangelogModelResize(), TestChangelogText(), TestNewChangelogModelViewContainsTitleAndNotes() (+12 more)

### Community 164 - "Candidate"
Cohesion: 0.11
Nodes (24): collectDiffAdditions(), pathCovered(), runDiffCatalog(), TestPathCovered(), TestUncoveredCandidates(), uncoveredCandidates(), selectAndEditConfig(), addManualDirectories() (+16 more)

### Community 166 - "dialLiveSession"
Cohesion: 0.24
Nodes (10): dialLiveSession(), applyEdit(), checkAndApply(), editAgain(), fetchConfig(), putConfig(), runEditConfig(), readPutContent() (+2 more)

### Community 168 - "mustSubkey"
Cohesion: 0.31
Nodes (9): newFileVaultAEAD(), openFileVault(), sealFileVault(), mustSubkey(), TestDeriveFileVaultSubkeyDeterministicAndSeparated(), TestIsFileVaultRecordShape(), TestOpenFileVaultRejectsTampering(), TestSealFileVaultNoncesNeverRepeat() (+1 more)

### Community 169 - "bpfstats/main.go"
Cohesion: 0.40
Nodes (8): chownToSudoUser(), compare(), dumpFullLog(), dumpTrace(), loadOne(), main(), measure(), short()

### Community 170 - "writeEvent"
Cohesion: 0.33
Nodes (8): UIDResolver, NewUIDResolver(), writeEvent(), allowedEvent(), deniedEvent(), TestWriteEventAllowed(), TestWriteEventAllowedSuppressed(), TestWriteEventBlocked()

### Community 172 - "Limitations"
Cohesion: 0.22
Nodes (9): Already-running processes, Binaries you can write, Event masks are not confinement, Interpreters and JITs, Limitations, Project status, Root, Side effects to expect with the daemon (+1 more)

### Community 173 - "Troubleshooting"
Cohesion: 0.22
Nodes (9): A program I use is denied, After a package update replaced a binary, An app breaks and nothing is logged, Backups, Check the guard by hand, Encryption key errors, Is the daemon running and enforcing?, Snapshots or disk backups stopped working (+1 more)

### Community 174 - "watchPattern"
Cohesion: 0.39
Nodes (6): addSystemBinary(), catalogPatterns(), catalogWatchPlan(), systemBinaryPatterns(), TestCatalogWatchPlanWatchesInspectors(), watchPattern

### Community 176 - "Daemon configuration"
Cohesion: 0.25
Nodes (8): Daemon configuration, Generic Electron apps: `[electron_apps]`, Library trust: `[libraries "<name>"]`, Per-binary event masks, Process inspectors, Reloading, Several guarded trees, one vault, Watch sections

### Community 177 - "edit-protected"
Cohesion: 0.25
Nodes (8): edit-protected, Editing the daemon config, Editor keys, Exit audit, Non-interactive writes, Offline and live mode, Temporary access for other programs: `--forward`, The password and the control socket

### Community 178 - "runLockdown"
Cohesion: 0.33
Nodes (7): configureDaemonLogging(), runLockdown(), dirIsRootOwnedSafe(), isBpffs(), mountBpffs(), ResolvePinBase(), TestDirIsRootOwnedSafe()

### Community 179 - "launchKey"
Cohesion: 0.36
Nodes (3): TrustGuard, launchKey(), TestLaunchRulesFitTheKey()

### Community 180 - "Compatibility"
Cohesion: 0.29
Nodes (7): Architectures, BPF-LSM must be active, not just compiled in, Check your host first, Compatibility, Distributions, Kernel floors per mode, Mandatory and best-effort hooks

### Community 181 - "How it works"
Cohesion: 0.25
Nodes (8): eBPF at three levels, Fail closed, Global flags, How it works, Identity is the executable's inode, Output, Watching from a browser: `--serve`, Where it sits: DAC, MAC and app-listener

### Community 182 - "parseConfig"
Cohesion: 0.29
Nodes (5): parseConfig(), BtrfsMounted(), mountinfoHasFstype(), TestMountinfoHasFstype(), progressReader

### Community 183 - "ParseEventsFlag"
Cohesion: 0.29
Nodes (6): ParseEventsFlag(), TestValidateFlagsEditConfig(), validateFlags(), TestValidateForwardFlags(), validateForwardFlags(), validateForwardRule()

### Community 184 - "bpfLsmListed"
Cohesion: 0.33
Nodes (7): bpfLsmListed(), CheckBPFLSM(), CheckBPFLSMAt(), TestBpfLsmListed(), TestBpfLsmListedMissingFile(), TestCheckBPFLSM(), writeLSM()

### Community 185 - "findTopLevelWindow"
Cohesion: 0.38
Nodes (4): findTopLevelWindow(), floatWindow(), internAtom(), setFloatHint()

### Community 186 - "Binary trust"
Cohesion: 0.33
Nodes (6): Binary ledger, Binary replacement, Binary trust, Launch scan, Multicall binaries, trust-binaries

### Community 187 - "guard: block file access by program"
Cohesion: 0.33
Nodes (6): Executables inside the guarded tree, Flags, guard: block file access by program, Multicall binaries, Replaced binaries, Whitelist and blacklist

### Community 188 - "network-guard: block network access by program"
Cohesion: 0.33
Nodes (6): Code integrity (`-w`), Flags, Multicall binaries, network-guard: block network access by program, Replaced binaries, What code integrity cannot stop

### Community 189 - "launch_rule_at"
Cohesion: 0.40
Nodes (3): launch_rule_at(), launch_step(), launch_token()

### Community 190 - "Development"
Cohesion: 0.40
Nodes (5): Demo GIFs, Development, Docker, Makefile targets, Tests

### Community 191 - "app-listener documentation"
Cohesion: 0.40
Nodes (5): app-listener documentation, Getting started, Modes, Operating, The daemon in depth

### Community 192 - "Self-updating apps and the catalog"
Cohesion: 0.50
Nodes (4): Catalog refresh, Self-updating apps, Self-updating apps and the catalog, The catalog

### Community 193 - "network-monitor: see what a program talks to"
Cohesion: 0.50
Nodes (4): Events, Flags, Multicall binaries, network-monitor: see what a program talks to

### Community 194 - "unlockUnderGuard"
Cohesion: 0.50
Nodes (5): lockVaultFully(), unlockUnderGuard(), vaultOpUnderGuard(), vaultScan(), SectionScan

### Community 195 - "CleanupStalePins"
Cohesion: 0.40
Nodes (5): CleanupStalePins(), genOfPinFile(), TestCleanupStalePins(), TestCleanupStalePinsMissingBase(), TestGenOfPinFile()

### Community 197 - "monitor: see what touches a path"
Cohesion: 0.67
Nodes (3): Events, Flags, monitor: see what touches a path

### Community 201 - "copyXattrs"
Cohesion: 0.50
Nodes (4): TestCopyXattrsPreservesUserAttributes(), copyXattrs(), listXattrNames(), splitXattrNames()

### Community 203 - "VerifyInstalledResourcesLocked"
Cohesion: 0.67
Nodes (3): TestVerifyResourcesLocked(), VerifyInstalledResourcesLocked(), VerifyResourcesLocked()

## Knowledge Gaps
- **142 isolated node(s):** `kernel-prof.sh script`, `tempJournalRow`, `launchRule`, `trustEvent`, `btrfsIoctlFsInfoArgs` (+137 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 487 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **42 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `absPath()` connect `absPath` to `daemon/daemon.go`, `TestMain`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`?**
  _High betweenness centrality (0.060) - this node is a cross-community bridge._
- **Why does `fileEditModel` connect `fileEditModel` to `daemon/daemon.go`, `fileedit_model_test.go`, `github.com/charmbracelet/bubbles/textarea.Model`, `codeedit.go`, `os.File`, `github.com/charmbracelet/bubbletea.Cmd`, `highlight.go`, `fileedit_symlink_test.go`?**
  _High betweenness centrality (0.029) - this node is a cross-community bridge._
- **Why does `IntegrationSuite` connect `IntegrationSuite` to `.newGuardTestContainer`?**
  _High betweenness centrality (0.021) - this node is a cross-community bridge._
- **What connects `kernel-prof.sh script`, `tempJournalRow`, `launchRule` to the rest of the system?**
  _142 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `daemon/daemon.go` be split into smaller, more focused modules?**
  _Cohesion score 0.09888201454466515 - nodes in this community are weakly interconnected._
- **Should `testing.T` be split into smaller, more focused modules?**
  _Cohesion score 0.023376937984496124 - nodes in this community are weakly interconnected._
- **Should `IntegrationSuite` be split into smaller, more focused modules?**
  _Cohesion score 0.0898838004101162 - nodes in this community are weakly interconnected._