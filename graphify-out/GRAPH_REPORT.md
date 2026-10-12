# Graph Report - app-listener  (2026-10-12)

## Corpus Check
- 349 files · ~1,552,204 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 12 file(s) not represented in the graph (top: (none) 7, .service 2, .pub 1)

## Summary
- 4569 nodes · 18055 edges · 201 communities (168 shown, 33 thin omitted)
- Extraction: 93% EXTRACTED · 7% INFERRED · 0% AMBIGUOUS · INFERRED: 1305 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `9107f5f6`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- daemon/daemon.go
- testing.T
- IntegrationSuite
- filevault_test.go
- daemonconfig_test.go
- networkguard.bpf.c
- monitor.bpf.c
- User
- guardUnitTest
- mkdirs
- auth.go
- planUpdaters
- TestMonitorUseCaseLifecycle
- IntegrationSuite
- applyDiffAdditions
- Vault
- IntegrationSuite
- guard_trust.bpf.c
- GenerateConf
- fileEditModel
- networkmonitor.bpf.c
- usecase/daemon_test.go
- Resource
- reloadOnce
- Monitor
- guard.bpf.c
- NetworkMonitor
- changelog.go
- SanitizeText
- buildOneGuard
- monitorUnitTest
- Serve
- GuardInodeKey
- IntegrationSuite
- fileedit_model_test.go
- DaemonEvent
- Config
- filevault.go
- controlServer
- bpfstats/main.go
- github.com/charmbracelet/bubbles/textarea.Model
- runOfflineEdit
- IntegrationSuite
- liveSession
- IntegrationSuite
- update_test.go
- highlight_test.go
- .syncSetsLocked
- TrustGuard
- tempRow
- inode_dev
- IntegrationSuite
- netModel
- github.com/charmbracelet/bubbletea.Cmd
- absPath
- codeedit.go
- sync.Mutex
- .SetGlobReservations
- SuperblockDev
- downloadFile
- netGuardModel
- exe_supersede.h
- fakeGuardRepo
- errno
- pinstate.go
- KernelDev
- time.Duration
- runNetworkMonitor
- catalogRefresher
- fcntl
- copyRegularAt
- Run
- integrationtests/guard_test.go
- btrfsLayoutFrom
- Ledger
- github.com/spf13/cobra.Command
- runMonitor
- process_vm_readv.c
- .step
- ListUsers
- sanitizeTerminalText
- CandidateDir
- MountVouch
- os.FileMode
- render.go
- watchSet
- TempBinary
- daemon mode (fscrypt + whitelist lifecycle)
- FORWARD temporary -w/-b whitelist grants
- LibraryClosure
- RefreshSection
- netinject.c
- NetGuard
- check-compatibility.sh
- BPF verifier 1M-insn complexity budget
- Inode-based (dev:ino) executable identity
- NewEventFanout
- NetGuard
- IntegrationSuite
- Well-known bypass classes
- Vault
- jit_provenance.c
- inspector_probe.c
- ptrace_race.c
- inodeChain
- Build & sign release assets (reusable workflow)
- RawTarget
- classifySupportError
- net_tester/main.go
- bottomBar
- app-listener Terminal Demo Recording
- update.go
- find_test.go
- configedit_test.go
- IntegrationSuite
- IntegrationSuite
- Regenerate eBPF bindings in pinned builder image
- auditGroupCoverage
- trustManager
- WebSocket /ws client
- NewDumpHook
- launch_probe.c
- install.sh
- mc_tag.h
- io.Writer
- make test-integration (rootful Docker suite)
- EventType
- supersede_probe.c
- bufio.Reader
- trace-app-libs.sh
- IntegrationSuite
- fakeDaemon
- Bun TMPDIR redirect with reserved .bun-* name
- Daemon self-protection (selfguards.go)
- github.com/Virgula0/app-listener
- network-monitor mode
- uninstall subcommand
- writeWithin
- RefreshCatalog
- os.File
- Guard
- secureResources
- launchKey
- SelectOrphans
- runGenKey
- NetworkGuardUseCase
- daemon: permanent protection with encryption at rest
- DumpHook
- Candidate
- releaseAssets
- .pathInfo
- loadDaemonConfig
- Installation
- ParseEventsFlag
- graphify-refresh.sh
- StripPinnedTempAllows
- .degradeWatchLoop
- tempGrant
- mountVouchesOwnership
- GuardEvent
- tempGrantSetup
- vouchedDevs
- ResolveLibraryClosure
- Vault
- sanityCheckBinary
- Limitations
- Troubleshooting
- github.com/stretchr/testify/suite.Suite
- runNetworkGuard
- Daemon configuration
- edit-protected
- runDaemon
- copyXattrs
- Compatibility
- How it works
- harnessSuite
- golang.org/x/sys/unix.Timespec
- bpfLsmListed
- newPathCache
- Binary trust
- guard: block file access by program
- network-guard: block network access by program
- VerifyInstalledResourcesLocked
- Development
- app-listener documentation
- Self-updating apps and the catalog
- network-monitor: see what a program talks to
- MulticallError
- monitor: see what touches a path
- kernel-prof.sh
- vaultFS
- find.go

## God Nodes (most connected - your core abstractions)
1. `Config` - 93 edges
2. `Guard` - 93 edges
3. `GuardInodeKey` - 89 edges
4. `Load()` - 85 edges
5. `absPath()` - 78 edges
6. `IntegrationSuite` - 75 edges
7. `IntegrationSuite` - 66 edges
8. `writeConfig()` - 65 edges
9. `fileEditModel` - 60 edges
10. `User` - 56 edges

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

## Communities (201 total, 33 thin omitted)

### Community 0 - "daemon/daemon.go"
Cohesion: 0.09
Nodes (20): Execute(), tempJournalRow, launchRule, trustEvent, btrfsIoctlFsInfoArgs, infraContainerPaths(), TestIntegrationSuite(), pooledContainer (+12 more)

### Community 1 - "testing.T"
Cohesion: 0.02
Nodes (141): newTestVetter(), testBinary(), TestVetterBootstrapRecordsThenChecks(), TestVetterLinkKeyedOnLink(), TestVetterRecordLiveUnderLinkLines(), TestVetterRefusesNewLineAfterBootstrap(), TestVetterResolverOnlyApproved(), TestVetterSwapWithinGenerationLosesRights() (+133 more)

### Community 2 - "IntegrationSuite"
Cohesion: 0.06
Nodes (13): IntegrationSuite, IntegrationSuite, configBinaryLines(), IntegrationSuite, parseDaemonEvents(), probeExeAs(), procGateMatrix(), procReadProbe() (+5 more)

### Community 3 - "filevault_test.go"
Cohesion: 0.18
Nodes (26): looksLikeFileVaultRecord(), newFileVaultAEAD(), openFileVault(), sealFileVault(), inode(), mustSubkey(), TestDeriveFileVaultSubkeyDeterministicAndSeparated(), TestIsEncryptedDispatchesToFileVaultForRegularFiles() (+18 more)

### Community 4 - "daemonconfig_test.go"
Cohesion: 0.08
Nodes (69): Load(), ResolvePendingPaths(), hasBinary(), hasPendingBinary(), TestElectronAppsBlockParsed(), TestElectronAppsBlockRejectsMalformedLines(), TestEncryptionGroupsSkipsLibDirs(), TestInspectorsBlockBounded() (+61 more)

### Community 5 - "networkguard.bpf.c"
Cohesion: 0.10
Nodes (48): admit_refusal(), check_watched(), code_vouched(), emit_event(), emit_gate(), env_loader(), env_step(), exe_inode_of() (+40 more)

### Community 6 - "monitor.bpf.c"
Cohesion: 0.14
Nodes (31): emit_event(), emit_event_kern(), fill_path(), read_path(), trace_do_faccessat(), trace_do_sendfile(), trace_do_splice(), trace_do_splice_direct() (+23 more)

### Community 7 - "User"
Cohesion: 0.10
Nodes (45): trustManager, askBunTmpdirUsers(), bunEntryConfigured(), bunLaunchersFor(), gatherPerUserSetup(), setupBunTmpdir(), deploy(), addKeysToAgent() (+37 more)

### Community 8 - "guardUnitTest"
Cohesion: 0.05
Nodes (22): addErrKind, guardUnitTest, ReplacementCheck, BinariesSummary(), classifyAddErr(), eventMask(), guardModeKey(), addRecorder() (+14 more)

### Community 9 - "mkdirs"
Cohesion: 0.13
Nodes (36): newWatchPattern(), newFakeInotify(), TestWatchSetBinaryWriteTriggersNotMkdir(), TestWatchSetFloodIsCapped(), TestWatchSetLibDirAppearanceTriggers(), TestWatchSetMissingRootWatchesDeepestAncestor(), TestWatchSetOverflowRescansAndRebuilds(), TestWatchSetSymlinkCreateTriggers() (+28 more)

### Community 10 - "auth.go"
Cohesion: 0.06
Nodes (53): Hash(), HashFileExists(), LoadHashFile(), OriginOf(), parseHash(), RemoveHashFile(), hashFileInode(), TestHashFileLifecycle() (+45 more)

### Community 11 - "planUpdaters"
Cohesion: 0.16
Nodes (16): hardLinked(), libBinaries(), mayUpdate(), planUpdaters(), resourceBinaries(), resourceLinks(), rules(), TestGeneralToolsToWarn_HardLinkIsNoUpdaterButNoWarning() (+8 more)

### Community 12 - "TestMonitorUseCaseLifecycle"
Cohesion: 0.24
Nodes (11): NewGuardUseCase(), NewNetworkMonitorUseCase(), newFakeMonitorRepo(), newFakeNetworkGuardRepo(), newFakeNetworkMonitorRepo(), TestGuardUseCaseLifecycle(), TestGuardUseCaseStartError(), TestMonitorUseCaseLifecycle() (+3 more)

### Community 13 - "IntegrationSuite"
Cohesion: 0.16
Nodes (7): netMonitorEvent, IntegrationSuite, fmtEvents(), IntegrationSuite, netMonitorDeltaEvents(), netMonitorHasTypes(), parseNetMonitorEvents()

### Community 14 - "applyDiffAdditions"
Cohesion: 0.05
Nodes (48): reloadDaemonForPasswordChange(), applyDiffAdditions(), restoreDaemonAfterDiffAbort(), promptEditPassword(), applyLiveRefresh(), buildBinaryIfNeeded(), checkRunningBinaryMatchesInstalled(), cleanupBackups() (+40 more)

### Community 15 - "Vault"
Cohesion: 0.12
Nodes (20): deprovisionKind, isFileVaultCiphertext(), isRegularFileTarget(), classifyDeprovision(), classifyDeprovisionErr(), Vault, hasEncryptionPolicy(), isLockedRegularFileErr() (+12 more)

### Community 16 - "IntegrationSuite"
Cohesion: 0.15
Nodes (3): IntegrationSuite, guardDeltaEvents(), le64HexKey()

### Community 17 - "guard_trust.bpf.c"
Cohesion: 0.09
Nodes (58): caller_updates(), current_exe_flags(), current_exe_inode(), current_exe_key(), current_writer_bits(), dentry_is_dir(), dentry_parent(), deny_if_protected() (+50 more)

### Community 18 - "GenerateConf"
Cohesion: 0.05
Nodes (55): TestDiffMergePreservesExistingAndParses(), discordCatalogEntry(), groupedDiscordConf(), TestAskEncryptionSkipsNeedEncryptionFalse(), TestAskSSHAgentUsersNoGuardedSSH(), TestCollectFilesystemPrereqsNoPanic(), TestCollectFilesystemPrereqsSkipsRegularFiles(), TestGeneratedLibraryBlockRoundTrips() (+47 more)

### Community 19 - "fileEditModel"
Cohesion: 0.10
Nodes (9): clipLabel(), humanSize(), clamp(), sortNodes(), chownUser, fileEditMode, fileEditModel, fileKind (+1 more)

### Community 20 - "networkmonitor.bpf.c"
Cohesion: 0.26
Nodes (16): BPF_KRETPROBE(), emit_event(), get_netns(), get_socket_proto(), is_watched_binary(), read_inet_addr(), trace_accept(), trace_accept4() (+8 more)

### Community 21 - "usecase/daemon_test.go"
Cohesion: 0.15
Nodes (47): NewDaemonUseCase(), partitionEncryptionRoots(), newFakeVault(), resource(), startDaemon(), TestDaemonUseCaseConstructorMismatch(), TestDaemonUseCaseEventsMerged(), TestDaemonUseCaseGrantEditAccess() (+39 more)

### Community 22 - "Resource"
Cohesion: 0.10
Nodes (9): Resource, GuardRepository, concurrencyLimit(), encryptionRootOf(), DaemonUseCase, uniqueEncryptionRoots(), unlockRoots(), vaultOpForGuard() (+1 more)

### Community 23 - "reloadOnce"
Cohesion: 0.11
Nodes (24): awaitStartupOrSignal(), backingDeviceUnion(), buildConcurrency(), buildGuards(), catchLifecycleSignals(), makeReloadHandler(), reloadOnce(), reloadSlotsNeeded() (+16 more)

### Community 25 - "guard.bpf.c"
Cohesion: 0.11
Nodes (57): add_inode_to_guard(), btrfs_copy_gate(), chmod_only_drops_write(), code_suspect(), count_degrade(), discover_guarded_parent(), emit_process_denial(), emit_with() (+49 more)

### Community 26 - "NetworkMonitor"
Cohesion: 0.11
Nodes (10): NetBpfEvent, NetEvent, FormatAddr(), Ntohs(), NewNetworkMonitor(), statInodeKey(), TestNetworkMonitor_NewFailure(), TestStatInodeKey() (+2 more)

### Community 27 - "changelog.go"
Cohesion: 0.17
Nodes (11): changelogText(), newChangelogModel(), showChangelog(), TestChangelogModelOtherKeysDoNotQuit(), TestChangelogModelQuitKeys(), TestChangelogModelResize(), TestChangelogText(), TestNewChangelogModelViewContainsTitleAndNotes() (+3 more)

### Community 28 - "SanitizeText"
Cohesion: 0.15
Nodes (8): multicallAdmissible(), newBinaryVetter(), openBinaryVetter(), ownerPaths(), binaryVetter, vetDecision, Source, SanitizeText()

### Community 29 - "buildOneGuard"
Cohesion: 0.07
Nodes (42): buildOneGuard(), eventFilterOptions(), systemPatterns(), systemRule(), TestEventFilterOptions(), newSelfGuards(), selfProtectSpecs(), resolveGuardConfig() (+34 more)

### Community 31 - "Serve"
Cohesion: 0.08
Nodes (28): basicAuth(), configureWebSocket(), decodeEndsAtEOF(), isInteractiveTerminal(), newServeListener(), newServeServer(), parseResize(), readResizeLoop() (+20 more)

### Community 32 - "GuardInodeKey"
Cohesion: 0.05
Nodes (32): inodeBatchFn, refusedReplacements, TestNearestRootClaims(), TestTieRank(), appendInodeRow(), batchInodeRows(), engine, inodeRowsOf() (+24 more)

### Community 34 - "fileedit_model_test.go"
Cohesion: 0.14
Nodes (28): OpenRegularNoFollow(), fileModeToUnixMode(), isBinaryContent(), isOctalMode(), newFileEditModel(), inodeOf(), maxLineWidth(), selectNamed() (+20 more)

### Community 35 - "DaemonEvent"
Cohesion: 0.10
Nodes (32): UIDResolver, NewUIDResolver(), drainEphemeral(), runHeadless(), writeEvent(), admitEvent(), foldable(), newGateLogLimiter() (+24 more)

### Community 36 - "Config"
Cohesion: 0.08
Nodes (48): sameConfig(), newDaemonModel(), runDaemonUI(), runServedTUI(), runTUI(), grantNewBinaries(), watchGroup, Line() (+40 more)

### Community 37 - "filevault.go"
Cohesion: 0.23
Nodes (14): classifyRegularFileTarget(), clearRecovery(), deriveFileVaultSubkey(), EnsureRecoverySidecarPlaceholder(), fileVaultRecoverPath(), Vault, openFileVaultWithMasterKey(), recoverFileInPlace() (+6 more)

### Community 38 - "controlServer"
Cohesion: 0.22
Nodes (5): controlServer, readClientLines(), streamEvents(), editControlSession, grantRequest

### Community 39 - "bpfstats/main.go"
Cohesion: 0.17
Nodes (18): guardSpec(), TestGuardSpecKeepsIoctlCompatOnlyWhereKernelHasIt(), verifyCollection(), VerifyLoad(), InertProgram(), KernelHasFunc(), TestKernelHasFunc(), trustSpec() (+10 more)

### Community 40 - "github.com/charmbracelet/bubbles/textarea.Model"
Cohesion: 0.26
Nodes (3): reposition(), Finder, overlay()

### Community 41 - "runOfflineEdit"
Cohesion: 0.12
Nodes (17): TestPickEncryptedDirSingleSkipsPicker(), pickEncryptedDir(), runOfflineEdit(), askDecrypt(), cleanOrphanedMetadata(), decryptDirectories(), decryptStep(), pickDirsToDecrypt() (+9 more)

### Community 44 - "IntegrationSuite"
Cohesion: 0.18
Nodes (7): guardBinaryFlag(), IntegrationSuite, IntegrationSuite, guardNetTypesForComm(), netGuardBlockedEventCount(), netGuardHasBlockedEvent(), netGuardTail()

### Community 45 - "update_test.go"
Cohesion: 0.11
Nodes (12): parseChecksum(), parsePublicKey(), rsaPublicKeyPEM(), signedChecksum(), testKeyPair(), TestNewerThanStable(), TestParseChecksum(), TestParsePublicKey() (+4 more)

### Community 46 - "highlight_test.go"
Cohesion: 0.16
Nodes (19): New(), TestCursorDoesNotBlink(), fg(), lexerFor(), NewHighlighter(), shebangLexer(), editSession(), sampleGo() (+11 more)

### Community 48 - ".syncSetsLocked"
Cohesion: 0.13
Nodes (16): resInfo, taintOwner, nextLabel(), TestDeliverLabelsCrossResourceProcessGates(), engine, isSubset(), planMemberRows(), planTaintOwners() (+8 more)

### Community 49 - "TrustGuard"
Cohesion: 0.09
Nodes (11): trustHook, trustRows, guardLSMHooks(), TestGuardLSMHooksIoctlCompatOnlyWhenPresent(), foreignPinnedLink(), setSupersedeTrust(), cStr(), TrustGuard (+3 more)

### Community 50 - "tempRow"
Cohesion: 0.24
Nodes (5): exeRow, tempMaskOp, tempRow, Guard, planTempBlock()

### Community 51 - "inode_dev"
Cohesion: 0.15
Nodes (18): exe_is_superseded(), exe_refused(), ctx_ptr(), inode_dev(), sb_dev(), check_and_emit_args(), current_is_inspector(), exe_tag_bits() (+10 more)

### Community 52 - "IntegrationSuite"
Cohesion: 0.06
Nodes (15): IntegrationSuite, decodeMulticall(), IntegrationSuite, multicallBypasses(), multicallConfig(), uutilsStrays(), vettedGroupIndex(), exploitTest (+7 more)

### Community 53 - "netModel"
Cohesion: 0.13
Nodes (16): Stats, parseKBValue(), parseStat(), parseStatus(), Read(), TestReadConcurrentRace(), TestStatsFieldsNonZero(), TestStatsValuesMonotonic() (+8 more)

### Community 54 - "github.com/charmbracelet/bubbletea.Cmd"
Cohesion: 0.10
Nodes (13): tickDaemonStats(), listenForGuardEvents(), tickGuardStats(), syncViewport(), renderBullets(), TestRenderBulletsVerbatimAndIndented(), guardEventLine, guardEventMsg (+5 more)

### Community 55 - "absPath"
Cohesion: 0.09
Nodes (8): IntegrationSuite, inNosuidNamespace(), IntegrationSuite, launchCase, absPath(), TestMain(), helperChild(), TestMain()

### Community 56 - "codeedit.go"
Cohesion: 0.19
Nodes (23): class, classOf(), column(), dedent(), Edit(), leadingIndent(), lines(), moveTo() (+15 more)

### Community 57 - "sync.Mutex"
Cohesion: 0.17
Nodes (6): controlManager, controlManager, newControlManager(), startControlManager(), startControlServer(), configEditor

### Community 58 - ".SetGlobReservations"
Cohesion: 0.16
Nodes (10): GlobChild, statFunc, compileGlobNames(), globText(), TrustGuard, resolveBits(), resolveChildren(), statKey() (+2 more)

### Community 59 - "SuperblockDev"
Cohesion: 0.20
Nodes (9): BtrfsMounted(), mountinfoDev(), mountinfoHasFstype(), OnBtrfs(), SuperblockDev(), SuperblockDevPath(), TestMountinfoDev(), TestMountinfoHasFstype() (+1 more)

### Community 60 - "downloadFile"
Cohesion: 0.22
Nodes (8): downloadFile(), downloadReleaseFiles(), fetchReleases(), TestDownloadFileHTTPError(), TestDownloadFileMode0700(), TestFetchReleases(), TestFetchReleasesHTTPError(), progressReader

### Community 61 - "netGuardModel"
Cohesion: 0.24
Nodes (5): NetGuardEvent, protoString(), netGuardEventLine, netGuardEventMsg, netGuardModel

### Community 62 - "exe_supersede.h"
Cohesion: 0.23
Nodes (14): exe_last_link_key(), exe_note_freed(), exe_note_supersede(), exe_seq_now(), exe_stamp_exec(), exe_stamp_fork(), exe_stamp_free(), exe_stamp_lost() (+6 more)

### Community 63 - "fakeGuardRepo"
Cohesion: 0.08
Nodes (4): Guard, SupersededBinary, supersededMark(), fakeGuardRepo

### Community 64 - "errno"
Cohesion: 0.21
Nodes (5): put(), wait_then_read(), dump_via_debugfs(), main(), probe_open()

### Community 65 - "pinstate.go"
Cohesion: 0.17
Nodes (23): ensurePinStateFilePlaceholder(), ensurePlaceholder(), pinOwnerLikelyAlive(), readPinState(), recoverPinState(), pinStateInode(), TestPinOwnerLikelyAlive(), TestPinOwnerLikelyAliveDeadPID() (+15 more)

### Community 66 - "KernelDev"
Cohesion: 0.13
Nodes (18): BinaryStat, binaryStatOf(), KernelDev(), LstatInode(), StatBinary(), BtrfsDevices(), btrfsSysfsDevices(), TestBtrfsSysfsDevices() (+10 more)

### Community 67 - "time.Duration"
Cohesion: 0.12
Nodes (11): newTestSession(), TestRunSessionIdleTimeoutResetsOnPing(), TestRunSessionStreamsEventsAsActivity(), debouncer, New(), recv(), TestSubDeliversMissedTickLate(), TestSubFiresEveryPeriod() (+3 more)

### Community 68 - "runNetworkMonitor"
Cohesion: 0.14
Nodes (8): runHeadless(), runHeadless(), runNetworkMonitor(), runTUI(), ProtocolString(), NetworkMonitorRepository, NewNetModel(), NetworkMonitorUseCase

### Community 69 - "catalogRefresher"
Cohesion: 0.21
Nodes (5): drain(), newCatalogRefresher(), readEvents(), swapFile(), catalogRefresher

### Community 70 - "fcntl"
Cohesion: 0.13
Nodes (4): denied(), main(), denied(), main()

### Community 71 - "copyRegularAt"
Cohesion: 0.31
Nodes (10): copyDirContents(), copyEntry(), copyRegularAt(), copySymlinkAt(), copyTreeRoot(), fchownFromInfo(), futimesFromInfo(), readDirNames() (+2 more)

### Community 72 - "Run"
Cohesion: 0.07
Nodes (14): columnState, dataRow, dataRowRenderer, guiModel, headerRenderer, headerWidget, findTopLevelWindow(), floatWindow() (+6 more)

### Community 73 - "integrationtests/guard_test.go"
Cohesion: 0.23
Nodes (6): attrOpCase, eventsForPath(), guardEventTypes(), parseGuardEvents(), guardEvent, guardExploitTest

### Community 74 - "btrfsLayoutFrom"
Cohesion: 0.15
Nodes (11): BtrfsLayout, fakeTypes, typeSource, btrfsLayoutFrom(), findMember(), memberOffset(), ResolveBtrfsLayout(), btrfsTypes() (+3 more)

### Community 75 - "Ledger"
Cohesion: 0.13
Nodes (8): EnsurePlaceholders(), Ledger, Pending, JournalPath(), Open(), TestEnsurePlaceholdersNoInstallDir(), TestLedgerKeepsJournalFile(), TestLedgerRecordLookupPending()

### Community 76 - "github.com/spf13/cobra.Command"
Cohesion: 0.15
Nodes (19): AddServeFlags(), credentialFlagState(), ServeConfig, isLoopbackHost(), ParseServeFlags(), requestedBoolFlag(), validateEnabledServe(), validateServeAddress() (+11 more)

### Community 78 - "runMonitor"
Cohesion: 0.16
Nodes (11): MakeDisplayPaths(), runHeadless(), runMonitor(), runTUI(), MonitorRepository, NewModel(), assertViewFits(), TestViewsFitReportedTerminalSize() (+3 more)

### Community 79 - "process_vm_readv.c"
Cohesion: 0.40
Nodes (6): dump_memory(), find_all_rw_regions(), find_heap(), find_victim_by_fd(), main(), spawn_less_pty()

### Community 80 - ".step"
Cohesion: 0.20
Nodes (11): placedWalk, confineBelow(), inUserHome(), GlobSystemPlaced(), placedDir(), placedMatches(), readlinkFd(), rootOwnedStat() (+3 more)

### Community 81 - "ListUsers"
Cohesion: 0.08
Nodes (35): deletePostBackups(), restoreBackups(), offerBackupCleanup(), runUninstall(), detectSSHAgentUnits(), isInstallerSSHAgentUnit(), removeKeyAndEmptyDir(), removeMasterKey() (+27 more)

### Community 82 - "sanitizeTerminalText"
Cohesion: 0.13
Nodes (11): sanitizeTerminalText(), sanitizeTerminalTexts(), clipToWidth(), formatResourceBar(), formatType(), listenForEvents(), tickStats(), TestSanitizeTerminalText() (+3 more)

### Community 83 - "CandidateDir"
Cohesion: 0.07
Nodes (24): entryWriters(), reserveChain(), symlinkTargetGlob(), underResource(), globBuilder, ParseGlobName(), TestParseGlobName(), confSafeMatch() (+16 more)

### Community 84 - "MountVouch"
Cohesion: 0.10
Nodes (15): ensureBtrfsLayout(), deleteResKeys(), dropSupersededLocked(), PruneSuperseded(), supersedeMaps(), supersedePruneLoop(), warnLostStamps(), syncU32Map() (+7 more)

### Community 85 - "os.FileMode"
Cohesion: 0.12
Nodes (24): TestDaemonUnitWithMetadataOutput(), daemonUnitWithMetadataOutput(), dupFD(), installFile(), installFileAs(), installServices(), installSSHAgent(), installUnit() (+16 more)

### Community 86 - "render.go"
Cohesion: 0.14
Nodes (13): lineStyles, viewportState, baseStyle(), Highlighter, gutterLabel(), lineLen(), markStyle(), repeatSpaces() (+5 more)

### Community 87 - "watchSet"
Cohesion: 0.10
Nodes (17): refreshAdmits(), TestRefreshAdmitsHomeMatchesAndRootPlacedNamesOnly(), addSystemBinary(), catalogPatterns(), catalogWatchPlan(), inHome(), isSymlink(), newWatchSet() (+9 more)

### Community 88 - "TempBinary"
Cohesion: 0.14
Nodes (15): TempBinary, TemporaryGrant, TempRule, applyTempGrants(), closeTempBinaries(), daemonUseCase, TemporaryAccess, logTempGrant() (+7 more)

### Community 89 - "daemon mode (fscrypt + whitelist lifecycle)"
Cohesion: 0.18
Nodes (10): Daemon catalog refresh (catalogrefresh.go / catalogwatch.go), exe_supersede.h superseded-key tracking, glob_plant reserved glob name POC, launch_probe risky-launch POC, preload_wait code-suspect POC, supersede_probe reused-inode POC, Daemon boot ordering (Type=notify, Before=user sessions), daemon mode (fscrypt + whitelist lifecycle) (+2 more)

### Community 90 - "FORWARD temporary -w/-b whitelist grants"
Cohesion: 0.21
Nodes (9): CONFIG post-AUTH daemon.conf editing (configedit.go), FORWARD temporary -w/-b whitelist grants, Daemon self-key (GUARD_ALLOW_ROOT) bypass, swap_benign / swap_reader in-place binary swap POC, edit-protected --edit-config, edit-protected subcommand, edit-protected live mode, edit-protected --forward (+1 more)

### Community 91 - "LibraryClosure"
Cohesion: 0.14
Nodes (14): ldConfParser, defaultLibDirs(), elfInterp(), fileExists(), isScript(), ldConfInclude(), ldSoConfDirs(), LdSoPreloadPaths() (+6 more)

### Community 92 - "RefreshSection"
Cohesion: 0.21
Nodes (12): RefreshOptions, catalogEntryMatch(), findCatalogRoot(), findCatalogWatchSubPath(), isInsidePath(), ResolveCatalogEntry(), admitNew(), keepExisting() (+4 more)

### Community 93 - "netinject.c"
Cohesion: 0.17
Nodes (10): attach_child(), await_exec(), do_connect(), errname(), main(), sleep_ms(), traceme(), do_connect() (+2 more)

### Community 95 - "NetGuard"
Cohesion: 0.12
Nodes (12): NetEventType, FirstOnBtrfs(), NetEventTypes(), ParseNetEventType(), reasonText(), eventsetSummary(), eventTypeKey(), NetGuard (+4 more)

### Community 96 - "check-compatibility.sh"
Cohesion: 0.38
Nodes (11): fail(), load_config(), lsm_bpf_instructions(), note(), pass(), read_config(), require_config(), section() (+3 more)

### Community 97 - "BPF verifier 1M-insn complexity budget"
Cohesion: 0.20
Nodes (8): PR vs base verifier cost comparison, Verifier gate workflow, bpf_loop() helper, tools/bpfstats verifier cost tool, static __noinline BPF-to-BPF helpers with args struct, rename-over-guarded-file bypass, xres_move cross-resource move POC, Kernel floor (5.8 / 5.10 / 5.17)

### Community 98 - "Inode-based (dev:ino) executable identity"
Cohesion: 0.18
Nodes (7): CheckBPFLSM preflight, internal/infrastructure shared eBPF kernel, Ports & adapters layering (cmd -> usecase -> repository -> engines), jit_provenance ledger (guard_jit_origin), daemon --check / --verifier-only preflight, guard mode (LSM whitelist/blacklist), network-guard mode

### Community 99 - "NewEventFanout"
Cohesion: 0.22
Nodes (6): NewEventFanout(), TestEventFanoutDropOldestUnderPressure(), TestEventFanoutDuplicatesInOrderAndCloses(), TestEventFanoutStopClosesOutputs(), EventFanout, EventFanout[T]

### Community 100 - "NetGuard"
Cohesion: 0.27
Nodes (6): NetGuard, judgePasses(), parseProcs(), TestJudgePasses(), TestParseProcs(), procRec

### Community 102 - "Well-known bypass classes"
Cohesion: 0.24
Nodes (9): btrfs copy-ioctl gate (file_ioctl / file_ioctl_compat), Well-known bypass classes, btrfs_search POC, io_uring POC, open_by_handle_at POC, process_vm_readv POC (remaining monitor gap), raw_block_device POC, statonly stat/statx metadata leak (+1 more)

### Community 103 - "Vault"
Cohesion: 0.20
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

### Community 107 - "inodeChain"
Cohesion: 0.23
Nodes (15): fdStat, inodeChain(), isProcFD(), mountinfoShowsSbRoot(), mountShowsSbRoot(), physicalParent(), statFD(), keyOf() (+7 more)

### Community 108 - "Build & sign release assets (reusable workflow)"
Cohesion: 0.31
Nodes (8): Build & sign release assets (reusable workflow), app-listener-gui best-effort build, Promote to stable release workflow, Release workflow (pre-release per push to main), Release promote flow, make check-compatibility, scripts/install.sh one-line installer, update self-update subcommand

### Community 109 - "RawTarget"
Cohesion: 0.27
Nodes (10): BuildEBPFTargets(), RawTarget, isSubDir(), pathWithinDir(), ResolveTargets(), validatePair(), ValidateTargets(), WarnIgnoredFlags() (+2 more)

### Community 110 - "classifySupportError"
Cohesion: 0.18
Nodes (10): classifySetupError(), classifySupportError(), TestClassifySetupErrorGeneric(), TestClassifySetupErrorNotSetup(), TestClassifySetupErrorNotSupported(), TestClassifySupportErrorEncryptionNotEnabled(), TestClassifySupportErrorEncryptionNotEnabledF2fs(), TestClassifySupportErrorGeneric() (+2 more)

### Community 111 - "net_tester/main.go"
Cohesion: 0.21
Nodes (7): connFD(), tcpClient(), tcpServer(), tcpServerDelayed(), udpRecvLoop(), udpSendLoop(), waitForMarker()

### Community 112 - "bottomBar"
Cohesion: 0.29
Nodes (5): TestProgressReader(), bottomBar, newBottomBar(), renderBar(), WithBottomBar()

### Community 113 - "app-listener Terminal Demo Recording"
Cohesion: 0.39
Nodes (8): BPF Build Pipeline (bpf2go + go build), CLI Subcommand Help (guard, install, monitor, network-guard, network-monitor), app-listener Terminal Demo Recording, Guard TUI Blacklist Mode (eBPF LSM), Guard TUI Whitelist Mode (eBPF LSM), File System Monitor TUI (eBPF), Startup Banner and eBPF Availability Check, TUI Resource Footer (Mem / CPU / Evt/s)

### Community 114 - "update.go"
Cohesion: 0.19
Nodes (19): applyUpdate(), compareStableVersions(), downloadAndVerify(), filterChannel(), latestTime(), newerThanInstalled(), newerThanStable(), parseStableVersion() (+11 more)

### Community 115 - "find_test.go"
Cohesion: 0.55
Nodes (10): NewFinder(), findEditor(), keys(), TestFinderLeavesSaveToHost(), TestFinderReadOnlyIgnoresReplace(), TestFinderReplace(), TestFinderReplaceWithSupersetTerminates(), TestFinderSearchAllMarksEveryMatch() (+2 more)

### Community 116 - "configedit_test.go"
Cohesion: 0.47
Nodes (9): newTestConfigEditor(), readFile(), TestConfigEditorApply(), TestConfigEditorRefusals(), TestConfigEditorRefusesNewMulticallLine(), TestConfigEditorRestoresOnFailedReload(), TestRunConfigEditOverTheSocket(), TestRunConfigEditRefusedDuringALiveSession() (+1 more)

### Community 119 - "Regenerate eBPF bindings in pinned builder image"
Cohesion: 0.33
Nodes (5): CI workflow (lint + test on non-draft PRs), CI lint job (make lint), CI test job (make test), golangci-lint v2 configuration, Issue #45 (clang 14 guard_path_rename over budget)

### Community 120 - "auditGroupCoverage"
Cohesion: 0.18
Nodes (13): auditAfterEdit(), auditAfterEditWithConfig(), auditEntry(), auditGroupCoverage(), auditTree(), coveredByWatch(), findResource(), groupWatchPaths() (+5 more)

### Community 121 - "trustManager"
Cohesion: 0.13
Nodes (20): inspectorPaths(), resolveInspectors(), TestInspectorPathsOnlyRootPlaced(), TestResolveInspectorsKeepsOnlyAdmitted(), buildTrustedSet(), TestBuildTrustedSetExcludesLibraryClosure(), closureRejections(), trustManager (+12 more)

### Community 122 - "WebSocket /ws client"
Cohesion: 0.29
Nodes (3): AppListener Shared Session page, WebSocket /ws client, --serve WebSocket mirroring

### Community 123 - "NewDumpHook"
Cohesion: 0.24
Nodes (11): applyVerbosity(), VerboseLevel, BridgeStdLog(), NewDumpHook(), TestBridgeStdLog(), TestDefaultVerboseMirrorsCurrentDisplay(), TestDumpHookFileIsSynced(), TestDumpHookFiltersByVerbosity() (+3 more)

### Community 124 - "launch_probe.c"
Cohesion: 0.53
Nodes (4): copy(), is_verb(), main(), usage()

### Community 125 - "install.sh"
Cohesion: 0.57
Nodes (6): die(), download(), info(), install.sh script, usage(), warn()

### Community 126 - "mc_tag.h"
Cohesion: 0.27
Nodes (9): mc_attest(), mc_basename(), mc_leader_start(), mc_stamp_fork(), mc_stamp_free(), mc_tag_bits(), netg_exec_applet(), netm_exec_applet() (+1 more)

### Community 127 - "io.Writer"
Cohesion: 0.27
Nodes (10): logRefreshChange(), TestLogRefreshChange(), candidates(), confirmable(), confirmAll(), describe(), hashNow(), run() (+2 more)

### Community 128 - "make test-integration (rootful Docker suite)"
Cohesion: 0.40
Nodes (4): make test-integration (rootful Docker suite), app-listener docker-compose service, ptrace_race POC, Docker usage (docker compose)

### Community 129 - "EventType"
Cohesion: 0.10
Nodes (13): BpfEvent, EventType, eventUnitTest, FileEvent, parseEvents(), TestAddBinaryActionsRefusesBeforeWriting(), TestCheckBinaryEventsReadOnlyRejectsRestriction(), Run() (+5 more)

### Community 130 - "supersede_probe.c"
Cohesion: 0.90
Nodes (4): answer(), main(), read_secret(), run_child()

### Community 131 - "bufio.Reader"
Cohesion: 0.18
Nodes (12): controlServer, readConfigPut(), TestReadGrantRequestConfig(), readAuthRequest(), readForwardRequest(), readGrantRequest(), readPathLines(), TestReadAuthRequest() (+4 more)

### Community 132 - "trace-app-libs.sh"
Cohesion: 0.60
Nodes (3): conf_directive_values(), list_conf_binaries(), trace-app-libs.sh script

### Community 135 - "Bun TMPDIR redirect with reserved .bun-* name"
Cohesion: 0.50
Nodes (4): lib_probe reserved library POC, Bun TMPDIR redirect with reserved .bun-* name, install --diff-catalog, install TUI wizard

### Community 142 - "writeWithin"
Cohesion: 0.22
Nodes (9): atomicWriteAt(), descendNoFollow(), refuseExistingSymlink(), TestWriteWithin_PlainFile(), TestWriteWithin_RefusesSymlinkDirComponent(), TestWriteWithin_RefusesSymlinkTarget(), TestWriteWithin_RejectsEscape(), writeAndSync() (+1 more)

### Community 143 - "RefreshCatalog"
Cohesion: 0.22
Nodes (9): updateCatalogConfig(), TestUnifiedDiff(), TestUnifiedDiffIdentical(), UnifiedDiff(), RefreshCatalog(), TestRefreshCatalogAdmitVetsOnlyNewLines(), TestRefreshCatalogKeepsSymlinkLineNotTarget(), TestRefreshCatalogReportsVersionBump() (+1 more)

### Community 144 - "os.File"
Cohesion: 0.05
Nodes (58): checkFreshDir(), fileID(), openBunTmp(), openDirNoFollow(), renewDir(), TestOpenBunTmp_ReplacesUnvettedDir(), TestOpenBunTmp_SymlinkReplacedTargetUntouched(), bunTmp (+50 more)

### Community 146 - "Guard"
Cohesion: 0.06
Nodes (21): binaryVerifyState, deferredBinary, rootHandle, sharedHashEntry, BinaryEntry, openRoot(), checkBtrfsKeys(), SharedPinDegraded() (+13 more)

### Community 147 - "secureResources"
Cohesion: 0.22
Nodes (10): encryptDirectories(), secureResources(), verifyEncryptionState(), collectFilesystemPrereqs(), confirmRunPrereq(), resolveFilesystemPrereqs(), askEncryption(), EnsureSystemSetup() (+2 more)

### Community 148 - "launchKey"
Cohesion: 0.36
Nodes (3): TrustGuard, launchKey(), TestLaunchRulesFitTheKey()

### Community 149 - "SelectOrphans"
Cohesion: 0.16
Nodes (20): modifiedContextWithSource(), appProtectorSet(), CleanOrphanedMetadata(), CleanOrphans(), isAppProtector(), listPolicies(), listProtectors(), LivePolicyDescriptors() (+12 more)

### Community 150 - "runGenKey"
Cohesion: 0.22
Nodes (10): confirmMasterKeyOverwrite(), runGenKey(), ensureInstalledBinary(), ensureMasterKey(), prepareInstallation(), TestEnsureInstalledBinary(), checkKeyLen(), GenerateMasterKey() (+2 more)

### Community 151 - "NetworkGuardUseCase"
Cohesion: 0.22
Nodes (4): NetworkGuardRepository, NetworkGuardUseCase, NewNetworkGuardUseCase(), TestNetworkGuardUseCaseStartError()

### Community 152 - "daemon: permanent protection with encryption at rest"
Cohesion: 0.17
Nodes (12): Boot ordering, daemon: permanent protection with encryption at rest, Filesystem-wide gates, Flags, Lifecycle, Process inspectors and taint, Quiet metadata denials, Related (+4 more)

### Community 153 - "DumpHook"
Cohesion: 0.32
Nodes (4): formatEntry(), sortedKeys(), sortStrings(), DumpHook

### Community 154 - "Candidate"
Cohesion: 0.10
Nodes (31): appendSectionsAndEdit(), collectDiffAdditions(), configPaths(), pathCovered(), runDiffCatalog(), TestPathCovered(), TestUncoveredCandidates(), uncoveredCandidates() (+23 more)

### Community 155 - "releaseAssets"
Cohesion: 0.32
Nodes (8): assetsFor(), checkUpdatePreconditions(), resolveAssets(), TestAssetsFor(), TestResolveAssetsPicksArch(), TestValidateChannel(), validateChannel(), releaseAssets

### Community 156 - ".pathInfo"
Cohesion: 0.39
Nodes (3): evalSymlinksOrEmpty(), pathCache, pathCacheEntry

### Community 157 - "loadDaemonConfig"
Cohesion: 0.29
Nodes (7): loadDaemonConfig(), setConfinedHomes(), TestLoadDaemonConfigEmptyResourcesIsCriticalStartup(), TestLoadDaemonConfigMissingFileIsCriticalStartup(), SetUserHomes(), TestResolveConfinedRootOwnedParents(), TestOpenSystemPlacedRefusesUserHome()

### Community 158 - "Installation"
Cohesion: 0.18
Nodes (11): Build from source, Bun-based apps (opencode), Docker, fscrypt prerequisites, Install flags, Installation, One-line installer, SSH agent (+3 more)

### Community 159 - "ParseEventsFlag"
Cohesion: 0.29
Nodes (6): ParseEventsFlag(), TestValidateFlagsEditConfig(), validateFlags(), TestValidateForwardFlags(), validateForwardFlags(), validateForwardRule()

### Community 161 - "StripPinnedTempAllows"
Cohesion: 0.38
Nodes (3): loadPinnedExeMaps(), stripPinnedAllow(), StripPinnedTempAllows()

### Community 164 - "mountVouchesOwnership"
Cohesion: 0.29
Nodes (7): hostView(), mountVouchesOwnership(), rootOwnedInode(), rootOwnedSafe(), SystemTrusted(), TestHostView(), TestMountVouchesOwnershipJudgesHostNamespace()

### Community 165 - "GuardEvent"
Cohesion: 0.07
Nodes (19): holdHeadless(), HeadlessLine(), runGuardHeadless(), bpfGuardEvent, commMatchesGuardedBinary(), decodeGuardEvent(), fsGateLabel(), GuardEvent (+11 more)

### Community 166 - "tempGrantSetup"
Cohesion: 0.57
Nodes (7): journalTo(), tempBinary(), tempGrantSetup(), TestGrantTemporaryAccessJournalsBeforeApplyAndClearsAfterRevoke(), TestGrantTemporaryAccessRefusals(), TestGrantTemporaryAccessRevokesOnInPlaceRewrite(), TestGrantTemporaryAccessRollsBackOnApplyFailure()

### Community 167 - "vouchedDevs"
Cohesion: 0.33
Nodes (6): isFuseType(), parseMajorMinor(), TestParseMajorMinorMatchesStatEncoding(), TestVouchedDevs(), TestVouchedDevsFailsClosed(), vouchedDevs()

### Community 168 - "ResolveLibraryClosure"
Cohesion: 0.29
Nodes (7): ResolveLibraryClosure(), skipWithoutHostRoot(), TestResolveLibraryClosure(), TestResolveLibraryClosureForeignArch(), TestResolveLibraryClosureScript(), TestResolveLibraryClosureStaticIsEmpty(), writeFileHelper()

### Community 170 - "sanityCheckBinary"
Cohesion: 0.40
Nodes (4): sanityCheckBinary(), TestSanityCheckBinaryRefusesOtherMachine(), runsNatively(), TestRunsNatively()

### Community 172 - "Limitations"
Cohesion: 0.22
Nodes (9): Already-running processes, Binaries you can write, Event masks are not confinement, Interpreters and JITs, Limitations, Project status, Root, Side effects to expect with the daemon (+1 more)

### Community 173 - "Troubleshooting"
Cohesion: 0.22
Nodes (9): A program I use is denied, After a package update replaced a binary, An app breaks and nothing is logged, Backups, Check the guard by hand, Encryption key errors, Is the daemon running and enforcing?, Snapshots or disk backups stopped working (+1 more)

### Community 175 - "runNetworkGuard"
Cohesion: 0.16
Nodes (16): ParseNetEventsFlag(), computeGuardBinaries(), confirmUnsafe(), discoverInfra(), resolveGuardMode(), runNetworkGuard(), runTUI(), DiscoverInfraBinaries() (+8 more)

### Community 176 - "Daemon configuration"
Cohesion: 0.25
Nodes (8): Daemon configuration, Generic Electron apps: `[electron_apps]`, Library trust: `[libraries "<name>"]`, Per-binary event masks, Process inspectors, Reloading, Several guarded trees, one vault, Watch sections

### Community 177 - "edit-protected"
Cohesion: 0.25
Nodes (8): edit-protected, Editing the daemon config, Editor keys, Exit audit, Non-interactive writes, Offline and live mode, Temporary access for other programs: `--forward`, The password and the control socket

### Community 178 - "runDaemon"
Cohesion: 0.08
Nodes (30): CheckBPFLSM(), CheckEBPF(), configureDaemonLogging(), lockOneRoot(), lockRootRecovering(), markReady(), newPinGeneration(), notifySystemdReady() (+22 more)

### Community 179 - "copyXattrs"
Cohesion: 0.50
Nodes (4): TestCopyXattrsPreservesUserAttributes(), copyXattrs(), listXattrNames(), splitXattrNames()

### Community 180 - "Compatibility"
Cohesion: 0.25
Nodes (8): Architectures, BPF-LSM must be active, not just compiled in, Check your host first, Compatibility, Distributions, Kernel floors per mode, Low-end hardware, Mandatory and best-effort hooks

### Community 181 - "How it works"
Cohesion: 0.25
Nodes (8): eBPF at three levels, Fail closed, Global flags, How it works, Identity is the executable's inode, Output, Watching from a browser: `--serve`, Where it sits: DAC, MAC and app-listener

### Community 184 - "bpfLsmListed"
Cohesion: 0.33
Nodes (7): bpfLsmListed(), CheckBPFLSM(), CheckBPFLSMAt(), TestBpfLsmListed(), TestBpfLsmListedMissingFile(), TestCheckBPFLSM(), writeLSM()

### Community 186 - "Binary trust"
Cohesion: 0.33
Nodes (6): Binary ledger, Binary replacement, Binary trust, Launch scan, Multicall binaries, trust-binaries

### Community 187 - "guard: block file access by program"
Cohesion: 0.33
Nodes (6): Executables inside the guarded tree, Flags, guard: block file access by program, Multicall binaries, Replaced binaries, Whitelist and blacklist

### Community 188 - "network-guard: block network access by program"
Cohesion: 0.33
Nodes (6): Code integrity (`-w`), Flags, Multicall binaries, network-guard: block network access by program, Replaced binaries, What code integrity cannot stop

### Community 189 - "VerifyInstalledResourcesLocked"
Cohesion: 0.67
Nodes (3): TestVerifyResourcesLocked(), VerifyInstalledResourcesLocked(), VerifyResourcesLocked()

### Community 190 - "Development"
Cohesion: 0.33
Nodes (6): Demo GIFs, Development, Docker, Makefile targets, Profiling, Tests

### Community 191 - "app-listener documentation"
Cohesion: 0.40
Nodes (5): app-listener documentation, Getting started, Modes, Operating, The daemon in depth

### Community 192 - "Self-updating apps and the catalog"
Cohesion: 0.50
Nodes (4): Catalog refresh, Self-updating apps, Self-updating apps and the catalog, The catalog

### Community 193 - "network-monitor: see what a program talks to"
Cohesion: 0.50
Nodes (4): Events, Flags, Multicall binaries, network-monitor: see what a program talks to

### Community 197 - "monitor: see what touches a path"
Cohesion: 0.67
Nodes (3): Events, Flags, monitor: see what touches a path

### Community 199 - "vaultFS"
Cohesion: 0.14
Nodes (14): TestRuntimeClassMarkers(), TestFilterExistingWhitelistSymlinkEscape(), AtomicWriteAt(), applyNewFileMeta(), defaultNewFileMode(), listEntries(), openRootT(), TestDefaultNewFileMode() (+6 more)

### Community 203 - "find.go"
Cohesion: 0.10
Nodes (11): findButton, span, diffModel, editorModel, findAll(), lowerRunes(), newFindInput(), runesEqual() (+3 more)

## Knowledge Gaps
- **144 isolated node(s):** `graphify-refresh.sh script`, `trustManager`, `controlManager`, `controlServer`, `tempJournalRow` (+139 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 478 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **33 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `absPath()` connect `absPath` to `daemon/daemon.go`, `IntegrationSuite`, `IntegrationSuite`, `integrationtests/guard_test.go`, `IntegrationSuite`, `auth.go`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`?**
  _High betweenness centrality (0.039) - this node is a cross-community bridge._
- **Why does `Guard` connect `Guard` to `daemon/daemon.go`, `GuardInodeKey`, `DaemonEvent`, `time.Duration`, `GuardEvent`, `fakeDaemon`, `os.File`, `reloadOnce`, `sync.Mutex`, `buildOneGuard`?**
  _High betweenness centrality (0.030) - this node is a cross-community bridge._
- **Why does `GuardInodeKey` connect `GuardInodeKey` to `GuardEvent`, `inodeChain`, `.syncSetsLocked`, `os.File`, `Guard`, `tempRow`, `TrustGuard`, `MountVouch`, `TempBinary`, `buildOneGuard`, `fakeGuardRepo`?**
  _High betweenness centrality (0.022) - this node is a cross-community bridge._
- **What connects `graphify-refresh.sh script`, `trustManager`, `controlManager` to the rest of the system?**
  _144 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `daemon/daemon.go` be split into smaller, more focused modules?**
  _Cohesion score 0.09344701335683775 - nodes in this community are weakly interconnected._
- **Should `testing.T` be split into smaller, more focused modules?**
  _Cohesion score 0.022413793103448276 - nodes in this community are weakly interconnected._
- **Should `IntegrationSuite` be split into smaller, more focused modules?**
  _Cohesion score 0.06158415841584158 - nodes in this community are weakly interconnected._