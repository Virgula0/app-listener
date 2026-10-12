# Graph Report - app-listener  (2026-10-12)

## Corpus Check
- 365 files · ~1,561,915 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 12 file(s) not represented in the graph (top: (none) 7, .service 2, .pub 1)

## Summary
- 4693 nodes · 19059 edges · 212 communities (178 shown, 34 thin omitted)
- Extraction: 93% EXTRACTED · 7% INFERRED · 0% AMBIGUOUS · INFERRED: 1364 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `b877b8d1`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- daemon/daemon.go
- testing.T
- IntegrationSuite
- SelectOrphans
- daemonconfig_test.go
- networkguard.bpf.c
- monitor.bpf.c
- rcUser
- guardUnitTest
- mkdirs
- User
- planUpdaters
- NetGuardEvent
- Candidate
- DeployInstalledBinary
- wipe
- IntegrationSuite
- guard_trust.bpf.c
- GenerateConf
- fileEditModel
- networkmonitor.bpf.c
- usecase/daemon_test.go
- Resource
- auth.go
- buildOneGuard
- guard.bpf.c
- NetworkMonitor
- find_test.go
- binaryVetter
- hashFileShared
- Monitor
- serve.go
- .AllowReplacement
- time.Duration
- fileedit_model_test.go
- IntegrationSuite
- Run
- runUninstall
- sync.Mutex
- bpfstats/main.go
- Finder
- filevault_test.go
- absPath
- liveSession
- IntegrationSuite
- update.go
- highlight_test.go
- .syncSetsLocked
- TrustGuard
- tempRow
- inode_dev
- main_test.go
- watchGroup
- github.com/charmbracelet/bubbletea.Cmd
- netGuardModel
- codeedit.go
- controlManager
- GuardInodeKey
- GuardEvent
- runNonInteractiveLivePut
- filevault.go
- exe_supersede.h
- fakeGuardRepo
- errno
- pinstate.go
- Vault
- globBuilder
- Config
- RunCmd
- fcntl
- copyRegularAt
- updateCatalogConfig
- time.Time
- btrfsLayoutFrom
- monitorUnitTest
- github.com/spf13/cobra.Command
- IntegrationSuite
- process_vm_readv.c
- golang.org/x/sys/unix.Stat_t
- downloadFile
- model
- newTestVetter
- MountVouch
- ResolvePinBase
- github.com/charmbracelet/bubbles/textarea.Model
- fakeDaemon
- mustSubkey
- daemon mode (fscrypt + whitelist lifecycle)
- FORWARD temporary -w/-b whitelist grants
- LibraryClosure
- DaemonEvent
- netinject.c
- NetGuard
- check-compatibility.sh
- BPF verifier 1M-insn complexity budget
- Inode-based (dev:ino) executable identity
- os.FileMode
- .judgeProcess
- IntegrationSuite
- Well-known bypass classes
- runNetworkGuard
- jit_provenance.c
- inspector_probe.c
- ptrace_race.c
- expandPlaceholders
- Build & sign release assets (reusable workflow)
- runMonitor
- fileedit_symlink_test.go
- net_tester/main.go
- NewEventFanout
- app-listener Terminal Demo Recording
- IntegrationSuite
- tempGrantSetup
- configedit_test.go
- IntegrationSuite
- IntegrationSuite
- Regenerate eBPF bindings in pinned builder image
- catalogRefresher
- NewDumpHook
- WebSocket /ws client
- netModel
- string
- install.sh
- mc_tag.h
- changelog.go
- make test-integration (rootful Docker suite)
- EventType
- supersede_probe.c
- runInstall
- trace-app-libs.sh
- IntegrationSuite
- serveWebSocket
- Bun TMPDIR redirect with reserved .bun-* name
- Daemon self-protection (selfguards.go)
- github.com/Virgula0/app-listener
- network-monitor mode
- uninstall subcommand
- TempBinary
- CandidateDir
- launchKey
- preload_wait.c
- Guard
- writeWithin
- runUpdate
- .SetGlobReservations
- runOfflineEdit
- StripPinnedTempAllows
- daemon: permanent protection with encryption at rest
- ConfinedEntry
- IntegrationSuite
- inodeChain
- highlight.go
- KernelDev
- Installation
- os.File
- graphify-refresh.sh
- classifySupportError
- ResolveLibraryClosure
- github.com/cilium/ebpf/link.Link
- BinaryStat
- IntegrationSuite
- is_event_type_allowed
- bottomBar
- runMaintenanceMode
- revertSSHAgents
- RefreshCatalog
- Limitations
- Troubleshooting
- ParseEventsFlag
- procRec
- Daemon configuration
- edit-protected
- IntegrationSuite
- Vault
- Compatibility
- How it works
- runForward
- loadDaemonConfig
- bpfLsmListed
- launch_rule_at
- Binary trust
- guard: block file access by program
- network-guard: block network access by program
- vouchedDevs
- Development
- app-listener documentation
- Self-updating apps and the catalog
- network-monitor: see what a program talks to
- parseResize
- formatEntry
- CleanupStalePins
- monitor: see what touches a path
- kernel-prof.sh
- vaultFS
- Vault
- HashCache
- harnessSuite
- .runEditorHarness
- copyXattrs
- parsePasswd
- golang.org/x/sys/unix.Timespec
- VerifyInstalledResourcesLocked
- MulticallError

## God Nodes (most connected - your core abstractions)
1. `Guard` - 95 edges
2. `Config` - 95 edges
3. `GuardInodeKey` - 91 edges
4. `Load()` - 89 edges
5. `IntegrationSuite` - 81 edges
6. `absPath()` - 79 edges
7. `writeConfig()` - 69 edges
8. `IntegrationSuite` - 66 edges
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

## Communities (212 total, 34 thin omitted)

### Community 0 - "daemon/daemon.go"
Cohesion: 0.10
Nodes (16): Execute(), tempJournalRow, launchRule, trustEvent, btrfsIoctlFsInfoArgs, Check(), collectDiagnostics(), readSysctl() (+8 more)

### Community 1 - "testing.T"
Cohesion: 0.02
Nodes (136): discordRefreshFixture(), TestRefreshOnceAppliesAndReloads(), TestRefreshOnceKeepsOperatorLinesAndSkipsNoOps(), TestRefreshOnceRestoresTheFileWhenTheReloadFails(), TestRefreshOnceSkipsAConfigChangedSinceLoad(), TestDebouncerCoalescesAndBoundsWait(), TestDebouncerMinGap(), TestEditControlSessionEndIdempotent() (+128 more)

### Community 2 - "IntegrationSuite"
Cohesion: 0.09
Nodes (10): configBinaryLines(), IntegrationSuite, parseDaemonEvents(), probeExeAs(), procGateMatrix(), procReadProbe(), rawExec(), shQuote() (+2 more)

### Community 3 - "SelectOrphans"
Cohesion: 0.20
Nodes (17): modifiedContextWithSource(), appProtectorSet(), CleanOrphans(), isAppProtector(), listPolicies(), listProtectors(), orphanedPolicies(), orphanedProtectors() (+9 more)

### Community 4 - "daemonconfig_test.go"
Cohesion: 0.08
Nodes (73): Load(), ResolvePendingPaths(), hasBinary(), hasPendingBinary(), TestElectronAppsBlockParsed(), TestElectronAppsBlockRejectsMalformedLines(), TestEncryptionGroupsSkipsLibDirs(), TestInspectorsBlockBounded() (+65 more)

### Community 5 - "networkguard.bpf.c"
Cohesion: 0.10
Nodes (38): admit_refusal(), check_watched(), code_vouched(), emit_gate(), env_loader(), env_step(), exe_inode_of(), fill_current_image() (+30 more)

### Community 6 - "monitor.bpf.c"
Cohesion: 0.15
Nodes (31): emit_event(), emit_event_kern(), fill_path(), read_path(), trace_do_faccessat(), trace_do_sendfile(), trace_do_splice(), trace_do_splice_direct() (+23 more)

### Community 7 - "rcUser"
Cohesion: 0.14
Nodes (32): rcTarget, bunTargets(), EnsureAddKeysToAgent(), EnsureBunLauncherEnv(), ensureRCBlock(), EnsureSSHAgentEnv(), openOwnedRegular(), rcBlock() (+24 more)

### Community 8 - "guardUnitTest"
Cohesion: 0.05
Nodes (25): addErrKind, guardUnitTest, ReplacementCheck, BinariesSummary(), classifyAddErr(), guardModeKey(), addRecorder(), copySelf() (+17 more)

### Community 9 - "mkdirs"
Cohesion: 0.13
Nodes (37): newWatchPattern(), newFakeInotify(), TestWatchSetBinaryWriteTriggersNotMkdir(), TestWatchSetFloodIsCapped(), TestWatchSetLibDirAppearanceTriggers(), TestWatchSetMissingRootWatchesDeepestAncestor(), TestWatchSetOverflowRescansAndRebuilds(), TestWatchSetSymlinkCreateTriggers() (+29 more)

### Community 10 - "User"
Cohesion: 0.13
Nodes (25): trustManager, askBunTmpdirUsers(), bunEntryConfigured(), bunLaunchersFor(), gatherPerUserSetup(), setupBunTmpdir(), deploy(), addKeysToAgent() (+17 more)

### Community 11 - "planUpdaters"
Cohesion: 0.16
Nodes (16): hardLinked(), libBinaries(), mayUpdate(), planUpdaters(), resourceBinaries(), resourceLinks(), rules(), TestGeneralToolsToWarn_HardLinkIsNoUpdaterButNoWarning() (+8 more)

### Community 12 - "NetGuardEvent"
Cohesion: 0.11
Nodes (17): NetGuardEvent, NetworkGuardRepository, NewGuardUseCase(), NetworkGuardUseCase, NewNetworkGuardUseCase(), NewNetworkMonitorUseCase(), newFakeMonitorRepo(), newFakeNetworkGuardRepo() (+9 more)

### Community 13 - "Candidate"
Cohesion: 0.13
Nodes (23): nextConfig(), appendSectionsAndEdit(), TestDiffMergePreservesExistingAndParses(), selectAndEditConfig(), addManualDirectories(), editConfig(), groupCandidates(), libraryBlocksFromCandidates() (+15 more)

### Community 14 - "DeployInstalledBinary"
Cohesion: 0.14
Nodes (17): buildBinaryIfNeeded(), installBinaryOnly(), installConfig(), mustCwd(), applyUpdate(), DaemonUnitInstalled(), DeployInstalledBinary(), DeployInstalledBinaryNoRestart() (+9 more)

### Community 15 - "wipe"
Cohesion: 0.12
Nodes (19): deprovisionKind, isFileVaultCiphertext(), isRegularFileTarget(), classifyDeprovision(), classifyDeprovisionErr(), hasEncryptionPolicy(), isLockedRegularFileErr(), isNotEncryptedErr() (+11 more)

### Community 16 - "IntegrationSuite"
Cohesion: 0.10
Nodes (9): attrOpCase, eventsForPath(), IntegrationSuite, guardDeltaEvents(), guardEventTypes(), le64HexKey(), parseGuardEvents(), guardEvent (+1 more)

### Community 17 - "guard_trust.bpf.c"
Cohesion: 0.11
Nodes (51): caller_updates(), current_exe_flags(), current_exe_inode(), current_exe_key(), current_writer_bits(), dentry_is_dir(), dentry_parent(), deny_if_protected() (+43 more)

### Community 18 - "GenerateConf"
Cohesion: 0.08
Nodes (37): TestSetSectionWhitelistPreservesGroupStructure(), LibraryBlock, Section, findSectionEnd(), findSectionStart(), firstContentLine(), GenerateConf(), isLibDirectiveText() (+29 more)

### Community 19 - "fileEditModel"
Cohesion: 0.12
Nodes (8): clipLabel(), humanSize(), sortNodes(), chownUser, fileEditMode, fileEditModel, fileKind, fileNode

### Community 20 - "networkmonitor.bpf.c"
Cohesion: 0.31
Nodes (16): BPF_KRETPROBE(), emit_event(), get_netns(), get_socket_proto(), is_watched_binary(), read_inet_addr(), trace_accept(), trace_accept4() (+8 more)

### Community 21 - "usecase/daemon_test.go"
Cohesion: 0.16
Nodes (45): NewDaemonUseCase(), newFakeVault(), resource(), startDaemon(), TestDaemonUseCaseConstructorMismatch(), TestDaemonUseCaseEventsMerged(), TestDaemonUseCaseGrantEditAccess(), TestDaemonUseCaseGroupedEncryptionRootsDeduplicated() (+37 more)

### Community 22 - "Resource"
Cohesion: 0.09
Nodes (10): Resource, GuardRepository, concurrencyLimit(), encryptionRootOf(), DaemonUseCase, partitionEncryptionRoots(), uniqueEncryptionRoots(), unlockRoots() (+2 more)

### Community 23 - "auth.go"
Cohesion: 0.12
Nodes (28): Hash(), HashFileExists(), LoadHashFile(), OriginOf(), parseHash(), RemoveHashFile(), hashFileInode(), TestHashFileLifecycle() (+20 more)

### Community 24 - "buildOneGuard"
Cohesion: 0.07
Nodes (39): buildOneGuard(), eventFilterOptions(), TestEventFilterOptions(), newSelfGuards(), selfProtectSpecs(), displayEntries(), holdForward(), holdHeadless() (+31 more)

### Community 25 - "guard.bpf.c"
Cohesion: 0.12
Nodes (51): add_inode_to_guard(), btrfs_copy_gate(), chmod_only_drops_write(), count_degrade(), discover_guarded_parent(), emit_process_denial(), evict_inode_from_guard(), get_dentry_from_path() (+43 more)

### Community 26 - "NetworkMonitor"
Cohesion: 0.07
Nodes (18): NetBpfEvent, NetEvent, NetEventType, Cstr(), FormatAddr(), NetEventTypes(), Ntohs(), ParseNetEventType() (+10 more)

### Community 27 - "find_test.go"
Cohesion: 0.55
Nodes (10): NewFinder(), findEditor(), keys(), TestFinderLeavesSaveToHost(), TestFinderReadOnlyIgnoresReplace(), TestFinderReplace(), TestFinderReplaceWithSupersetTerminates(), TestFinderSearchAllMarksEveryMatch() (+2 more)

### Community 28 - "binaryVetter"
Cohesion: 0.06
Nodes (28): newBinaryVetter(), openBinaryVetter(), newHashCache(), ds(), TestHashCacheFlushKeepsOnlyThisRunsRows(), candidates(), confirmable(), confirmAll() (+20 more)

### Community 29 - "hashFileShared"
Cohesion: 0.14
Nodes (11): binaryVerifyState, HashStore, sharedHashEntry, confinedEntry(), claimHash(), hashFileShared(), hashPathShared(), publishHash() (+3 more)

### Community 30 - "Monitor"
Cohesion: 0.11
Nodes (5): evalSymlinksOrEmpty(), NewMonitor(), Monitor, pathCache, pathCacheEntry

### Community 31 - "serve.go"
Cohesion: 0.15
Nodes (18): basicAuth(), isInteractiveTerminal(), newServeListener(), newServeServer(), relayServeSignals(), runServe(), securityHeaders(), Serve() (+10 more)

### Community 32 - ".AllowReplacement"
Cohesion: 0.18
Nodes (6): TrustGuard, UpdaterPlan, OnBtrfs(), SuperblockDev(), SuperblockDevPath(), TestSuperblockDevMatchesStatOffBtrfs()

### Community 33 - "time.Duration"
Cohesion: 0.13
Nodes (9): IntegrationSuite, pooledContainer, New(), recv(), TestSubDeliversMissedTickLate(), TestSubFiresEveryPeriod(), TestSubscribeDoesNotFireImmediately(), Pulse (+1 more)

### Community 34 - "fileedit_model_test.go"
Cohesion: 0.13
Nodes (30): OpenRegularNoFollow(), fileModeToUnixMode(), isBinaryContent(), isOctalMode(), newFileEditModel(), inodeOf(), maxLineWidth(), selectNamed() (+22 more)

### Community 35 - "IntegrationSuite"
Cohesion: 0.16
Nodes (7): netMonitorEvent, IntegrationSuite, fmtEvents(), IntegrationSuite, netMonitorDeltaEvents(), netMonitorHasTypes(), parseNetMonitorEvents()

### Community 36 - "Run"
Cohesion: 0.07
Nodes (14): columnState, dataRow, dataRowRenderer, guiModel, headerRenderer, headerWidget, findTopLevelWindow(), floatWindow() (+6 more)

### Community 37 - "runUninstall"
Cohesion: 0.17
Nodes (20): deletePostBackups(), restoreBackups(), offerBackupCleanup(), runUninstall(), removeKeyAndEmptyDir(), removeMasterKey(), TestRemoveKeyAndEmptyDir(), Delete() (+12 more)

### Community 38 - "sync.Mutex"
Cohesion: 0.11
Nodes (19): controlManager, controlServer, readConfigPut(), TestReadGrantRequestConfig(), controlServer, readAuthRequest(), readClientLines(), readForwardRequest() (+11 more)

### Community 39 - "bpfstats/main.go"
Cohesion: 0.15
Nodes (21): guardSpec(), TestGuardSpecKeepsIoctlCompatOnlyWhereKernelHasIt(), verifyCollection(), VerifyLoad(), applyInodeCapacity(), TestApplyInodeCapacity(), InertProgram(), KernelHasFunc() (+13 more)

### Community 40 - "Finder"
Cohesion: 0.15
Nodes (9): findButton, span, findAll(), Finder, lowerRunes(), newFindInput(), overlay(), runesEqual() (+1 more)

### Community 41 - "filevault_test.go"
Cohesion: 0.31
Nodes (17): looksLikeFileVaultRecord(), inode(), TestIsEncryptedDispatchesToFileVaultForRegularFiles(), TestIsProvisionedForFile(), TestLockUnlockFileInPlaceIdempotent(), TestStageRecoverySealsContentAndStaysInPlace(), TestUnlockFileInPlaceRecoversFromInterruptedTransform(), TestUnlockFileInPlaceWrongKeyFails() (+9 more)

### Community 42 - "absPath"
Cohesion: 0.08
Nodes (4): IntegrationSuite, inNosuidNamespace(), IntegrationSuite, absPath()

### Community 44 - "IntegrationSuite"
Cohesion: 0.18
Nodes (7): guardBinaryFlag(), IntegrationSuite, IntegrationSuite, guardNetTypesForComm(), netGuardBlockedEventCount(), netGuardHasBlockedEvent(), netGuardTail()

### Community 45 - "update.go"
Cohesion: 0.10
Nodes (21): compareStableVersions(), newerThanStable(), parseChecksum(), parsePublicKey(), parseStableVersion(), parseVerPart(), readInstalledVersion(), rsaPublicKeyPEM() (+13 more)

### Community 46 - "highlight_test.go"
Cohesion: 0.19
Nodes (23): moveTo(), New(), SetText(), TestBackspaceAndDeleteRemoveAnIndentStep(), TestCursorDoesNotBlink(), TestEditKeys(), TestEnterKeepsIndentation(), TestEnterPastNinetyNineLines() (+15 more)

### Community 48 - ".syncSetsLocked"
Cohesion: 0.13
Nodes (16): resInfo, taintOwner, nextLabel(), TestDeliverLabelsCrossResourceProcessGates(), engine, isSubset(), planMemberRows(), planTaintOwners() (+8 more)

### Community 49 - "TrustGuard"
Cohesion: 0.11
Nodes (8): trustHook, trustRows, guardLSMHooks(), TestGuardLSMHooksIoctlCompatOnlyWhenPresent(), setSupersedeTrust(), cStr(), TrustGuard, logTrustDenied()

### Community 50 - "tempRow"
Cohesion: 0.17
Nodes (7): exeRow, tempGrant, tempMaskOp, tempRow, Guard, liveTrust(), planTempBlock()

### Community 51 - "inode_dev"
Cohesion: 0.12
Nodes (23): exe_is_superseded(), exe_refused(), ctx_ptr(), inode_dev(), sb_dev(), check_and_emit_args(), code_suspect(), current_is_inspector() (+15 more)

### Community 52 - "main_test.go"
Cohesion: 0.05
Nodes (20): IntegrationSuite, decodeMulticall(), IntegrationSuite, multicallBypasses(), multicallConfig(), uutilsStrays(), vettedGroupIndex(), exploitTest (+12 more)

### Community 53 - "watchGroup"
Cohesion: 0.08
Nodes (44): grantNewBinaries(), watchGroup, Lines(), TestLinesKeyOnLinkAndIncludePending(), addResource(), applyAllowLib(), applyConfigLine(), applyDirective() (+36 more)

### Community 54 - "github.com/charmbracelet/bubbletea.Cmd"
Cohesion: 0.12
Nodes (7): clamp(), renderBullets(), TestRenderBulletsVerbatimAndIndented(), noticeModel, servedModel, serveTestModel, sessionModel

### Community 55 - "netGuardModel"
Cohesion: 0.14
Nodes (12): computeGuardBinaries(), Stats, parseKBValue(), parseStat(), parseStatus(), Read(), TestReadConcurrentRace(), TestStatsFieldsNonZero() (+4 more)

### Community 56 - "codeedit.go"
Cohesion: 0.25
Nodes (17): class, classOf(), column(), dedent(), Edit(), leadingIndent(), lines(), Navigate() (+9 more)

### Community 57 - "controlManager"
Cohesion: 0.29
Nodes (3): controlManager, newControlManager(), startControlManager()

### Community 58 - "GuardInodeKey"
Cohesion: 0.05
Nodes (26): refusedReplacements, TestTieRank(), appendInodeRow(), batchInodeRows(), engine, inodeRowsOf(), pathWithin(), TestPathWithinMatchesPrefixRule() (+18 more)

### Community 59 - "GuardEvent"
Cohesion: 0.07
Nodes (22): newTestSession(), TestRunSessionIdleTimeoutResetsOnPing(), TestRunSessionStreamsEventsAsActivity(), bpfGuardEvent, decodeGuardEvent(), fsGateLabel(), GuardEvent, logBacklogDenial() (+14 more)

### Community 60 - "runNonInteractiveLivePut"
Cohesion: 0.20
Nodes (10): auditAfterEdit(), auditEntry(), auditTree(), coveredByWatch(), relOrPath(), TestAuditTree(), TestWithin(), within() (+2 more)

### Community 61 - "filevault.go"
Cohesion: 0.23
Nodes (14): classifyRegularFileTarget(), clearRecovery(), deriveFileVaultSubkey(), EnsureRecoverySidecarPlaceholder(), fileVaultRecoverPath(), Vault, openFileVaultWithMasterKey(), recoverFileInPlace() (+6 more)

### Community 62 - "exe_supersede.h"
Cohesion: 0.17
Nodes (18): exe_last_link_key(), exe_note_freed(), exe_note_supersede(), exe_seq_now(), exe_stamp_exec(), exe_stamp_fork(), exe_stamp_free(), exe_stamp_lost() (+10 more)

### Community 64 - "errno"
Cohesion: 0.28
Nodes (3): dump_via_debugfs(), main(), probe_open()

### Community 65 - "pinstate.go"
Cohesion: 0.15
Nodes (25): configureDaemonLogging(), runLockdown(), ensurePinStateFilePlaceholder(), ensurePlaceholder(), pinOwnerLikelyAlive(), readPinState(), recoverPinState(), pinStateInode() (+17 more)

### Community 66 - "Vault"
Cohesion: 0.18
Nodes (10): migrationTarget, writeTempWithMetadata(), classifyMigrationTarget(), Vault, readClassifiedFile(), stampCopyMetadata(), copyTreeRoot(), CopyTreeWithProgress() (+2 more)

### Community 67 - "globBuilder"
Cohesion: 0.20
Nodes (8): entryWriters(), reserveChain(), underResource(), globBuilder, ParseGlobName(), TestParseGlobName(), inOwnTree(), TestCatalogTrustGlobsAreReservable()

### Community 68 - "Config"
Cohesion: 0.07
Nodes (46): setConfinedHomes(), inspectorPaths(), resolveInspectors(), TestInspectorPathsOnlyRootPlaced(), TestResolveInspectorsKeepsOnlyAdmitted(), buildTrustedSet(), TestBuildTrustedSetExcludesLibraryClosure(), closureRejections() (+38 more)

### Community 69 - "RunCmd"
Cohesion: 0.14
Nodes (18): reloadDaemonForPasswordChange(), preflightDeployedBinary(), revertSystemFiles(), FindDaemonProcesses(), isDaemonCmdline(), TestFindDaemonProcessesEmpty(), TestIsDaemonCmdline(), EnableAndVerify() (+10 more)

### Community 70 - "fcntl"
Cohesion: 0.16
Nodes (4): denied(), main(), denied(), main()

### Community 71 - "copyRegularAt"
Cohesion: 0.36
Nodes (9): copyDirContents(), copyEntry(), copyRegularAt(), copySymlinkAt(), fchownFromInfo(), futimesFromInfo(), readDirNames(), readlinkFD() (+1 more)

### Community 72 - "updateCatalogConfig"
Cohesion: 0.16
Nodes (16): dialLiveSession(), applyEdit(), checkAndApply(), editAgain(), fetchConfig(), putConfig(), runEditConfig(), pingInterval() (+8 more)

### Community 73 - "time.Time"
Cohesion: 0.09
Nodes (23): addSystemBinary(), catalogPatterns(), catalogWatchPlan(), isSymlink(), newWatchSet(), systemBinaryPatterns(), fakeInotify, watchDir (+15 more)

### Community 74 - "btrfsLayoutFrom"
Cohesion: 0.14
Nodes (12): BtrfsLayout, fakeTypes, typeSource, ensureBtrfsLayout(), btrfsLayoutFrom(), findMember(), memberOffset(), ResolveBtrfsLayout() (+4 more)

### Community 75 - "monitorUnitTest"
Cohesion: 0.10
Nodes (3): newPathCache(), dirTarget(), monitorUnitTest

### Community 76 - "github.com/spf13/cobra.Command"
Cohesion: 0.15
Nodes (19): AddServeFlags(), credentialFlagState(), ServeConfig, isLoopbackHost(), ParseServeFlags(), requestedBoolFlag(), validateEnabledServe(), validateServeAddress() (+11 more)

### Community 79 - "process_vm_readv.c"
Cohesion: 0.40
Nodes (6): dump_memory(), find_all_rw_regions(), find_heap(), find_victim_by_fd(), main(), spawn_less_pty()

### Community 80 - "golang.org/x/sys/unix.Stat_t"
Cohesion: 0.20
Nodes (11): placedWalk, confineBelow(), inUserHome(), GlobSystemPlaced(), placedDir(), placedMatches(), readlinkFd(), rootOwnedStat() (+3 more)

### Community 81 - "downloadFile"
Cohesion: 0.17
Nodes (9): downloadFile(), TestDownloadFileHTTPError(), TestDownloadFileMode0700(), BtrfsMounted(), mountinfoDev(), mountinfoHasFstype(), TestMountinfoDev(), TestMountinfoHasFstype() (+1 more)

### Community 82 - "model"
Cohesion: 0.22
Nodes (7): formatType(), listenForEvents(), syncViewport(), tickStats(), eventLine, eventMsg, model

### Community 83 - "newTestVetter"
Cohesion: 0.38
Nodes (10): newTestVetter(), testBinary(), TestVetterBootstrapRecordsThenChecks(), TestVetterLinkKeyedOnLink(), TestVetterRecordLiveUnderLinkLines(), TestVetterRefusesNewLineAfterBootstrap(), TestVetterResolverOnlyApproved(), TestVetterSeesInPlaceRewriteThroughSharedHashes() (+2 more)

### Community 84 - "MountVouch"
Cohesion: 0.08
Nodes (17): deleteInoKeys(), dropSupersededLocked(), exeGone(), Guard, PruneSuperseded(), supersededMark(), supersedeMaps(), supersedePruneLoop() (+9 more)

### Community 85 - "ResolvePinBase"
Cohesion: 0.50
Nodes (5): dirIsRootOwnedSafe(), isBpffs(), mountBpffs(), ResolvePinBase(), TestDirIsRootOwnedSafe()

### Community 86 - "github.com/charmbracelet/bubbles/textarea.Model"
Cohesion: 0.14
Nodes (14): lineStyles, viewportState, fg(), baseStyle(), Highlighter, gutterLabel(), lineLen(), markStyle() (+6 more)

### Community 88 - "mustSubkey"
Cohesion: 0.31
Nodes (9): newFileVaultAEAD(), openFileVault(), sealFileVault(), mustSubkey(), TestDeriveFileVaultSubkeyDeterministicAndSeparated(), TestIsFileVaultRecordShape(), TestOpenFileVaultRejectsTampering(), TestSealFileVaultNoncesNeverRepeat() (+1 more)

### Community 89 - "daemon mode (fscrypt + whitelist lifecycle)"
Cohesion: 0.18
Nodes (10): Daemon catalog refresh (catalogrefresh.go / catalogwatch.go), exe_supersede.h superseded-key tracking, glob_plant reserved glob name POC, launch_probe risky-launch POC, preload_wait code-suspect POC, supersede_probe reused-inode POC, Daemon boot ordering (Type=notify, Before=user sessions), daemon mode (fscrypt + whitelist lifecycle) (+2 more)

### Community 90 - "FORWARD temporary -w/-b whitelist grants"
Cohesion: 0.21
Nodes (9): CONFIG post-AUTH daemon.conf editing (configedit.go), FORWARD temporary -w/-b whitelist grants, Daemon self-key (GUARD_ALLOW_ROOT) bypass, swap_benign / swap_reader in-place binary swap POC, edit-protected --edit-config, edit-protected subcommand, edit-protected live mode, edit-protected --forward (+1 more)

### Community 91 - "LibraryClosure"
Cohesion: 0.10
Nodes (19): ldConfParser, defaultLibDirs(), elfInterp(), fileExists(), isScript(), ldConfInclude(), ldSoConfDirs(), LdSoPreloadPaths() (+11 more)

### Community 92 - "DaemonEvent"
Cohesion: 0.08
Nodes (38): UIDResolver, NewUIDResolver(), catchLifecycleSignals(), drainEphemeral(), newDaemonModel(), runDaemonUI(), runHeadless(), runServedTUI() (+30 more)

### Community 93 - "netinject.c"
Cohesion: 0.17
Nodes (10): attach_child(), await_exec(), do_connect(), errname(), main(), sleep_ms(), traceme(), do_connect() (+2 more)

### Community 95 - "NetGuard"
Cohesion: 0.18
Nodes (4): FirstOnBtrfs(), NetGuard, mcNameKey(), statInodeKey()

### Community 96 - "check-compatibility.sh"
Cohesion: 0.38
Nodes (11): fail(), load_config(), lsm_bpf_instructions(), note(), pass(), read_config(), require_config(), section() (+3 more)

### Community 97 - "BPF verifier 1M-insn complexity budget"
Cohesion: 0.20
Nodes (8): PR vs base verifier cost comparison, Verifier gate workflow, bpf_loop() helper, tools/bpfstats verifier cost tool, static __noinline BPF-to-BPF helpers with args struct, rename-over-guarded-file bypass, xres_move cross-resource move POC, Kernel floor (5.8 / 5.10 / 5.17)

### Community 98 - "Inode-based (dev:ino) executable identity"
Cohesion: 0.18
Nodes (7): CheckBPFLSM preflight, internal/infrastructure shared eBPF kernel, Ports & adapters layering (cmd -> usecase -> repository -> engines), jit_provenance ledger (guard_jit_origin), daemon --check / --verifier-only preflight, guard mode (LSM whitelist/blacklist), network-guard mode

### Community 99 - "os.FileMode"
Cohesion: 0.13
Nodes (22): TestDaemonUnitWithMetadataOutput(), daemonUnitWithMetadataOutput(), dupFD(), installFile(), installFileAs(), installServices(), installSSHAgent(), installSSHAgentEnv() (+14 more)

### Community 100 - ".judgeProcess"
Cohesion: 0.20
Nodes (11): reasonText(), codeVouched(), NetGuard, isTraced(), loaderEnv(), roFor(), TestCodeVouched(), TestLoaderEnv() (+3 more)

### Community 102 - "Well-known bypass classes"
Cohesion: 0.24
Nodes (9): btrfs copy-ioctl gate (file_ioctl / file_ioctl_compat), Well-known bypass classes, btrfs_search POC, io_uring POC, open_by_handle_at POC, process_vm_readv POC (remaining monitor gap), raw_block_device POC, statonly stat/statx metadata leak (+1 more)

### Community 103 - "runNetworkGuard"
Cohesion: 0.10
Nodes (26): CheckBPFLSM(), CheckEBPF(), ParseNetEventsFlag(), runBPFCheck(), confirmUnsafe(), discoverInfra(), resolveGuardMode(), runHeadless() (+18 more)

### Community 104 - "jit_provenance.c"
Cohesion: 0.47
Nodes (8): copy_into(), main(), make_temp_copy(), mode_hold(), mode_memfd(), mode_passfd(), mode_self(), report_dlopen()

### Community 105 - "inspector_probe.c"
Cohesion: 0.56
Nodes (8): check_attach(), check_exe(), check_mem(), check_root(), check_vmread(), main(), report(), traced_exec()

### Community 106 - "ptrace_race.c"
Cohesion: 0.31
Nodes (5): main(), peek_buffer(), run_tracer(), run_victim(), wait_for_file()

### Community 107 - "expandPlaceholders"
Cohesion: 0.19
Nodes (8): symlinkTargetGlob(), expandPlaceholders(), CandidateDir, CandidateDir, LibDirGlob, TrustGlob, reservedGlobs(), splitTrustGlob()

### Community 108 - "Build & sign release assets (reusable workflow)"
Cohesion: 0.31
Nodes (8): Build & sign release assets (reusable workflow), app-listener-gui best-effort build, Promote to stable release workflow, Release workflow (pre-release per push to main), Release promote flow, make check-compatibility, scripts/install.sh one-line installer, update self-update subcommand

### Community 109 - "runMonitor"
Cohesion: 0.11
Nodes (19): BuildEBPFTargets(), RawTarget, isSubDir(), MakeDisplayPaths(), pathWithinDir(), ResolveTargets(), validatePair(), ValidateTargets() (+11 more)

### Community 110 - "fileedit_symlink_test.go"
Cohesion: 0.40
Nodes (10): mustRead(), mustWrite(), swappedParent(), TestChmodRefusesSymlinkedParent(), TestCreateEntryRefusesPlantedSymlink(), TestCreateEntryRefusesSymlinkedParent(), TestDeleteRefusesSymlinkedParent(), TestSaveFollowsPinnedVaultNotSwappedAncestor() (+2 more)

### Community 111 - "net_tester/main.go"
Cohesion: 0.21
Nodes (7): connFD(), tcpClient(), tcpServer(), tcpServerDelayed(), udpRecvLoop(), udpSendLoop(), waitForMarker()

### Community 112 - "NewEventFanout"
Cohesion: 0.22
Nodes (6): NewEventFanout(), TestEventFanoutDropOldestUnderPressure(), TestEventFanoutDuplicatesInOrderAndCloses(), TestEventFanoutStopClosesOutputs(), EventFanout, EventFanout[T]

### Community 113 - "app-listener Terminal Demo Recording"
Cohesion: 0.39
Nodes (8): BPF Build Pipeline (bpf2go + go build), CLI Subcommand Help (guard, install, monitor, network-guard, network-monitor), app-listener Terminal Demo Recording, Guard TUI Blacklist Mode (eBPF LSM), Guard TUI Whitelist Mode (eBPF LSM), File System Monitor TUI (eBPF), Startup Banner and eBPF Availability Check, TUI Resource Footer (Mem / CPU / Evt/s)

### Community 115 - "tempGrantSetup"
Cohesion: 0.57
Nodes (7): journalTo(), tempBinary(), tempGrantSetup(), TestGrantTemporaryAccessJournalsBeforeApplyAndClearsAfterRevoke(), TestGrantTemporaryAccessRefusals(), TestGrantTemporaryAccessRevokesOnInPlaceRewrite(), TestGrantTemporaryAccessRollsBackOnApplyFailure()

### Community 116 - "configedit_test.go"
Cohesion: 0.47
Nodes (9): newTestConfigEditor(), readFile(), TestConfigEditorApply(), TestConfigEditorRefusals(), TestConfigEditorRefusesNewMulticallLine(), TestConfigEditorRestoresOnFailedReload(), TestRunConfigEditOverTheSocket(), TestRunConfigEditRefusedDuringALiveSession() (+1 more)

### Community 119 - "Regenerate eBPF bindings in pinned builder image"
Cohesion: 0.33
Nodes (5): CI workflow (lint + test on non-draft PRs), CI lint job (make lint), CI test job (make test), golangci-lint v2 configuration, Issue #45 (clang 14 guard_path_rename over budget)

### Community 120 - "catalogRefresher"
Cohesion: 0.13
Nodes (7): drain(), newCatalogRefresher(), readEvents(), sameConfig(), swapFile(), catalogRefresher, debouncer

### Community 121 - "NewDumpHook"
Cohesion: 0.20
Nodes (12): applyVerbosity(), VerboseLevel, BridgeStdLog(), NewDumpHook(), TestBridgeStdLog(), TestDefaultVerboseMirrorsCurrentDisplay(), TestDumpHookFileIsSynced(), TestDumpHookFiltersByVerbosity() (+4 more)

### Community 122 - "WebSocket /ws client"
Cohesion: 0.29
Nodes (3): AppListener Shared Session page, WebSocket /ws client, --serve WebSocket mirroring

### Community 123 - "netModel"
Cohesion: 0.12
Nodes (13): sanitizeTerminalText(), sanitizeTerminalTexts(), clipToWidth(), formatResourceBar(), formatAddr(), formatNetProto(), formatNetType(), listenForNetEvents() (+5 more)

### Community 124 - "string"
Cohesion: 0.16
Nodes (6): main(), report(), copy(), is_verb(), main(), usage()

### Community 125 - "install.sh"
Cohesion: 0.57
Nodes (6): die(), download(), info(), install.sh script, usage(), warn()

### Community 126 - "mc_tag.h"
Cohesion: 0.24
Nodes (10): mc_attest(), mc_basename(), mc_leader_start(), mc_stamp_fork(), mc_stamp_free(), mc_tag_bits(), netg_exec_applet(), netg_task_free() (+2 more)

### Community 127 - "changelog.go"
Cohesion: 0.22
Nodes (7): newChangelogModel(), showChangelog(), TestChangelogModelOtherKeysDoNotQuit(), TestChangelogModelQuitKeys(), TestChangelogModelResize(), TestNewChangelogModelViewContainsTitleAndNotes(), changelogModel

### Community 128 - "make test-integration (rootful Docker suite)"
Cohesion: 0.40
Nodes (4): make test-integration (rootful Docker suite), app-listener docker-compose service, ptrace_race POC, Docker usage (docker compose)

### Community 129 - "EventType"
Cohesion: 0.08
Nodes (13): BpfEvent, EventType, eventUnitTest, FileEvent, parseEvents(), TestAddBinaryActionsRefusesBeforeWriting(), TestCheckBinaryEventsReadOnlyRejectsRestriction(), Run() (+5 more)

### Community 130 - "supersede_probe.c"
Cohesion: 0.90
Nodes (4): answer(), main(), read_secret(), run_child()

### Community 131 - "runInstall"
Cohesion: 0.14
Nodes (15): confirmMasterKeyOverwrite(), runGenKey(), runOneShotMode(), checkRunningBinaryMatchesInstalled(), cleanupBackups(), ensureInstalledBinary(), ensureMasterKey(), prepareInstallation() (+7 more)

### Community 132 - "trace-app-libs.sh"
Cohesion: 0.60
Nodes (3): conf_directive_values(), list_conf_binaries(), trace-app-libs.sh script

### Community 134 - "serveWebSocket"
Cohesion: 0.22
Nodes (8): configureWebSocket(), readResizeLoop(), sameOrigin(), serveWebSocket(), TestSameOrigin(), writeSnapshot(), writeSnapshotLoop(), snapshotHub

### Community 135 - "Bun TMPDIR redirect with reserved .bun-* name"
Cohesion: 0.50
Nodes (4): lib_probe reserved library POC, Bun TMPDIR redirect with reserved .bun-* name, install --diff-catalog, install TUI wizard

### Community 142 - "TempBinary"
Cohesion: 0.10
Nodes (18): forwardRequest, TempBinary, TemporaryGrant, TempRule, applyTempGrants(), closeTempBinaries(), daemonUseCase, TemporaryAccess (+10 more)

### Community 143 - "CandidateDir"
Cohesion: 0.14
Nodes (12): systemPatterns(), systemRule(), confSafeMatch(), DiscoverSystem(), BinaryRule, CandidateDir, homeMatchConfined(), symlinkStaysInParent() (+4 more)

### Community 144 - "launchKey"
Cohesion: 0.36
Nodes (3): TrustGuard, launchKey(), TestLaunchRulesFitTheKey()

### Community 146 - "Guard"
Cohesion: 0.05
Nodes (26): lockVaultFully(), unlockUnderGuard(), vaultOpUnderGuard(), vaultScan(), deferredBinary, rootHandle, BinaryEntry, SectionScan (+18 more)

### Community 147 - "writeWithin"
Cohesion: 0.22
Nodes (9): atomicWriteAt(), descendNoFollow(), refuseExistingSymlink(), TestWriteWithin_PlainFile(), TestWriteWithin_RefusesSymlinkDirComponent(), TestWriteWithin_RefusesSymlinkTarget(), TestWriteWithin_RejectsEscape(), writeAndSync() (+1 more)

### Community 148 - "runUpdate"
Cohesion: 0.11
Nodes (27): changelogText(), TestChangelogText(), assetsFor(), checkUpdatePreconditions(), confirmUpdate(), downloadAndVerify(), downloadReleaseFiles(), fetchReleases() (+19 more)

### Community 149 - ".SetGlobReservations"
Cohesion: 0.20
Nodes (9): GlobChild, statFunc, compileGlobNames(), globText(), TrustGuard, resolveBits(), resolveChildren(), statKey() (+1 more)

### Community 150 - "runOfflineEdit"
Cohesion: 0.12
Nodes (16): TestPickEncryptedDirSingleSkipsPicker(), pickEncryptedDir(), runOfflineEdit(), askDecrypt(), decryptDirectories(), decryptStep(), pickDirsToDecrypt(), TestDecryptDirectoriesEmpty() (+8 more)

### Community 151 - "StripPinnedTempAllows"
Cohesion: 0.38
Nodes (3): loadPinnedExeMaps(), stripPinnedAllow(), StripPinnedTempAllows()

### Community 152 - "daemon: permanent protection with encryption at rest"
Cohesion: 0.17
Nodes (12): Boot ordering, daemon: permanent protection with encryption at rest, Filesystem-wide gates, Flags, Lifecycle, Process inspectors and taint, Quiet metadata denials, Related (+4 more)

### Community 153 - "ConfinedEntry"
Cohesion: 0.22
Nodes (22): ConfinedEntry(), HashStats(), newVerifyState(), hashedReused(), resetSharedHashes(), TestConcurrentCallersHashOnce(), TestConfinedEntryHashesOncePerFingerprint(), TestNewVerifyStateSeededOnlyByMatchingHash() (+14 more)

### Community 155 - "inodeChain"
Cohesion: 0.29
Nodes (12): inodeChain(), isProcFD(), mountinfoShowsSbRoot(), mountShowsSbRoot(), keyOf(), mkTree(), TestInodeChainDirIncludesItself(), TestInodeChainFollow() (+4 more)

### Community 156 - "highlight.go"
Cohesion: 0.13
Nodes (6): diffModel, editorModel, Highlighter, lexerFor(), shebangLexer(), tokenClassOf()

### Community 157 - "KernelDev"
Cohesion: 0.21
Nodes (11): fdStat, physicalParent(), statFD(), binaryStatOf(), KernelDev(), LstatInode(), BtrfsDevices(), btrfsSysfsDevices() (+3 more)

### Community 158 - "Installation"
Cohesion: 0.18
Nodes (11): Build from source, Bun-based apps (opencode), Docker, fscrypt prerequisites, Install flags, Installation, One-line installer, SSH agent (+3 more)

### Community 159 - "os.File"
Cohesion: 0.04
Nodes (69): multicallAdmissible(), ownerPaths(), checkFreshDir(), fileID(), openBunTmp(), openDirNoFollow(), renewDir(), plantBun() (+61 more)

### Community 161 - "classifySupportError"
Cohesion: 0.18
Nodes (10): classifySetupError(), classifySupportError(), TestClassifySetupErrorGeneric(), TestClassifySetupErrorNotSetup(), TestClassifySetupErrorNotSupported(), TestClassifySupportErrorEncryptionNotEnabled(), TestClassifySupportErrorEncryptionNotEnabledF2fs(), TestClassifySupportErrorGeneric() (+2 more)

### Community 162 - "ResolveLibraryClosure"
Cohesion: 0.18
Nodes (11): hostView(), mountVouchesOwnership(), ResolveLibraryClosure(), skipWithoutHostRoot(), TestHostView(), TestMountVouchesOwnershipJudgesHostNamespace(), TestResolveLibraryClosure(), TestResolveLibraryClosureForeignArch() (+3 more)

### Community 163 - "github.com/cilium/ebpf/link.Link"
Cohesion: 0.33
Nodes (3): foreignPinnedLink(), NetGuard, hook

### Community 164 - "BinaryStat"
Cohesion: 0.29
Nodes (10): BinaryStat, binaryStatOf(), StatBinary(), StatBinaryFile(), durableFS(), fsidOf(), hostVouches(), StatDurable() (+2 more)

### Community 166 - "is_event_type_allowed"
Cohesion: 0.44
Nodes (10): emit_event(), guard_net_socket_bind(), guard_net_socket_connect(), guard_net_socket_listen(), guard_net_socket_recvmsg(), guard_net_socket_sendmsg(), is_event_type_allowed(), is_socket_guarded() (+2 more)

### Community 167 - "bottomBar"
Cohesion: 0.29
Nodes (5): TestProgressReader(), bottomBar, newBottomBar(), renderBar(), WithBottomBar()

### Community 168 - "runMaintenanceMode"
Cohesion: 0.22
Nodes (10): applyLiveRefresh(), parseMaintenanceFlags(), runMaintenanceMode(), runUpdateCatalogOnly(), softenAutomatedRefreshErr(), TestValidateMaintenanceFlags(), validateMaintenanceFlags(), TestApplyLiveRefreshDeliversReload() (+2 more)

### Community 169 - "revertSSHAgents"
Cohesion: 0.25
Nodes (9): detectSSHAgentUnits(), isInstallerSSHAgentUnit(), removeSSHAgentEnv(), removeSSHAgentUnit(), revertSSHAgents(), staleRCUsers(), TestIsInstallerSSHAgentUnit(), TestRemoveSSHAgentUnit() (+1 more)

### Community 170 - "RefreshCatalog"
Cohesion: 0.25
Nodes (8): logRefreshChange(), TestLogRefreshChange(), SectionChange, RefreshCatalog(), TestRefreshCatalogAdmitVetsOnlyNewLines(), TestRefreshCatalogKeepsSymlinkLineNotTarget(), TestRefreshCatalogReportsVersionBump(), TestRefreshCatalogUserSectionsOnly()

### Community 172 - "Limitations"
Cohesion: 0.22
Nodes (9): Already-running processes, Binaries you can write, Event masks are not confinement, Interpreters and JITs, Limitations, Project status, Root, Side effects to expect with the daemon (+1 more)

### Community 173 - "Troubleshooting"
Cohesion: 0.22
Nodes (9): A program I use is denied, After a package update replaced a binary, An app breaks and nothing is logged, Backups, Check the guard by hand, Encryption key errors, Is the daemon running and enforcing?, Snapshots or disk backups stopped working (+1 more)

### Community 174 - "ParseEventsFlag"
Cohesion: 0.29
Nodes (6): ParseEventsFlag(), TestValidateFlagsEditConfig(), validateFlags(), TestValidateForwardFlags(), validateForwardFlags(), validateForwardRule()

### Community 175 - "procRec"
Cohesion: 0.33
Nodes (5): judgePasses(), parseProcs(), TestJudgePasses(), TestParseProcs(), procRec

### Community 176 - "Daemon configuration"
Cohesion: 0.22
Nodes (9): Daemon configuration, Daemon settings: `[settings]`, Generic Electron apps: `[electron_apps]`, Library trust: `[libraries "<name>"]`, Per-binary event masks, Process inspectors, Reloading, Several guarded trees, one vault (+1 more)

### Community 177 - "edit-protected"
Cohesion: 0.25
Nodes (8): edit-protected, Editing the daemon config, Editor keys, Exit audit, Non-interactive writes, Offline and live mode, Temporary access for other programs: `--forward`, The password and the control socket

### Community 179 - "Vault"
Cohesion: 0.08
Nodes (39): backingDeviceUnion(), buildConcurrency(), buildGuards(), lockOneRoot(), lockRootRecovering(), makeReloadHandler(), newPinGeneration(), prepareDaemonStart() (+31 more)

### Community 180 - "Compatibility"
Cohesion: 0.25
Nodes (8): Architectures, BPF-LSM must be active, not just compiled in, Check your host first, Compatibility, Distributions, Kernel floors per mode, Low-end hardware, Mandatory and best-effort hooks

### Community 181 - "How it works"
Cohesion: 0.25
Nodes (8): eBPF at three levels, Fail closed, Global flags, How it works, Identity is the executable's inode, Output, Watching from a browser: `--serve`, Where it sits: DAC, MAC and app-listener

### Community 182 - "runForward"
Cohesion: 0.40
Nodes (6): chooseResources(), confirmUserPlaced(), eventsLabel(), forwardBinaries(), normalizedEvents(), runForward()

### Community 183 - "loadDaemonConfig"
Cohesion: 0.17
Nodes (11): loadDaemonConfig(), markReady(), noteReloadedSettings(), notifySystemdReady(), resolveConfigPath(), TestLoadDaemonConfigEmptyResourcesIsCriticalStartup(), TestLoadDaemonConfigMissingFileIsCriticalStartup(), engine (+3 more)

### Community 184 - "bpfLsmListed"
Cohesion: 0.33
Nodes (7): bpfLsmListed(), CheckBPFLSM(), CheckBPFLSMAt(), TestBpfLsmListed(), TestBpfLsmListedMissingFile(), TestCheckBPFLSM(), writeLSM()

### Community 185 - "launch_rule_at"
Cohesion: 0.40
Nodes (3): launch_rule_at(), launch_step(), launch_token()

### Community 186 - "Binary trust"
Cohesion: 0.29
Nodes (7): Binary ledger, Binary replacement, Binary trust, Hash cache, Launch scan, Multicall binaries, trust-binaries

### Community 187 - "guard: block file access by program"
Cohesion: 0.33
Nodes (6): Executables inside the guarded tree, Flags, guard: block file access by program, Multicall binaries, Replaced binaries, Whitelist and blacklist

### Community 188 - "network-guard: block network access by program"
Cohesion: 0.33
Nodes (6): Code integrity (`-w`), Flags, Multicall binaries, network-guard: block network access by program, Replaced binaries, What code integrity cannot stop

### Community 189 - "vouchedDevs"
Cohesion: 0.33
Nodes (6): isFuseType(), parseMajorMinor(), TestParseMajorMinorMatchesStatEncoding(), TestVouchedDevs(), TestVouchedDevsFailsClosed(), vouchedDevs()

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

### Community 194 - "parseResize"
Cohesion: 0.40
Nodes (3): decodeEndsAtEOF(), parseResize(), TestParseResize()

### Community 195 - "formatEntry"
Cohesion: 0.50
Nodes (3): formatEntry(), sortedKeys(), sortStrings()

### Community 196 - "CleanupStalePins"
Cohesion: 0.40
Nodes (5): CleanupStalePins(), genOfPinFile(), TestCleanupStalePins(), TestCleanupStalePinsMissingBase(), TestGenOfPinFile()

### Community 197 - "monitor: see what touches a path"
Cohesion: 0.67
Nodes (3): Events, Flags, monitor: see what touches a path

### Community 199 - "vaultFS"
Cohesion: 0.14
Nodes (13): TestRuntimeClassMarkers(), TestFilterExistingWhitelistSymlinkEscape(), AtomicWriteAt(), applyNewFileMeta(), defaultNewFileMode(), listEntries(), openRootT(), TestDefaultNewFileMode() (+5 more)

### Community 204 - "copyXattrs"
Cohesion: 0.50
Nodes (4): TestCopyXattrsPreservesUserAttributes(), copyXattrs(), listXattrNames(), splitXattrNames()

### Community 205 - "parsePasswd"
Cohesion: 0.50
Nodes (4): parsePasswd(), shell(), TestParsePasswd(), TestParsePasswdMalformed()

### Community 207 - "VerifyInstalledResourcesLocked"
Cohesion: 0.67
Nodes (3): TestVerifyResourcesLocked(), VerifyInstalledResourcesLocked(), VerifyResourcesLocked()

## Knowledge Gaps
- **147 isolated node(s):** `tempJournalRow`, `launchRule`, `trustEvent`, `btrfsIoctlFsInfoArgs`, `attrOpCase` (+142 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 487 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **34 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `absPath()` connect `absPath` to `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `main_test.go`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`?**
  _High betweenness centrality (0.058) - this node is a cross-community bridge._
- **Why does `Guard` connect `Guard` to `daemon/daemon.go`, `time.Duration`, `sync.Mutex`, `time.Time`, `Vault`, `fakeDaemon`, `buildOneGuard`, `controlManager`, `GuardInodeKey`, `GuardEvent`, `DaemonEvent`, `hashFileShared`, `os.File`?**
  _High betweenness centrality (0.034) - this node is a cross-community bridge._
- **Why does `guardDeltaEvents()` connect `IntegrationSuite` to `absPath`?**
  _High betweenness centrality (0.026) - this node is a cross-community bridge._
- **What connects `tempJournalRow`, `launchRule`, `trustEvent` to the rest of the system?**
  _147 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `daemon/daemon.go` be split into smaller, more focused modules?**
  _Cohesion score 0.10321038714489118 - nodes in this community are weakly interconnected._
- **Should `testing.T` be split into smaller, more focused modules?**
  _Cohesion score 0.02219938335046249 - nodes in this community are weakly interconnected._
- **Should `IntegrationSuite` be split into smaller, more focused modules?**
  _Cohesion score 0.0898838004101162 - nodes in this community are weakly interconnected._