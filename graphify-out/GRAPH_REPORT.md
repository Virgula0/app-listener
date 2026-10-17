# Graph Report - app-listener  (2026-10-17)

## Corpus Check
- 435 files · ~1,614,033 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 13 file(s) not represented in the graph (top: (none) 8, .service 2, .pub 1)

## Summary
- 5496 nodes · 23356 edges · 203 communities (171 shown, 32 thin omitted)
- Extraction: 91% EXTRACTED · 9% INFERRED · 0% AMBIGUOUS · INFERRED: 2027 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `e73a3931`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- daemon/daemon.go
- testing.T
- IntegrationSuite
- GenerateConf
- daemonconfig_test.go
- networkguard.bpf.c
- monitor.bpf.c
- rcUser
- guardUnitTest
- mkdirs
- SanitizeText
- buildTrustedSet
- UI
- procstats.go
- os.FileMode
- SelectOrphans
- IntegrationSuite
- guard_trust.bpf.c
- IntegrationSuite
- fileEditModel
- networkmonitor.bpf.c
- usecase/daemon_test.go
- Resource
- auth.go
- NetGuardEvent
- guard.bpf.c
- Ledger
- networkPage
- inodeChain
- github.com/charmbracelet/bubbles/textarea.Model
- journal
- serve.go
- bpfstats/main.go
- update.go
- filevault.go
- github.com/spf13/cobra.Command
- IntegrationSuite
- BinaryEntry
- controlServer
- fileedit_model_test.go
- Guard
- IntegrationSuite
- absPath
- monitorUnitTest
- IntegrationSuite
- update_test.go
- daemon_multicall_test.go
- stdio
- golang.org/x/sys/unix.Stat_t
- sync.Mutex
- LibraryClosure
- runOfflineEdit
- .FilterExistingWhitelist
- Config
- github.com/charmbracelet/bubbletea.Cmd
- TrustGuard
- .showSized
- systemd.go
- GuardInodeKey
- IntegrationSuite
- buildOneGuard
- icon
- exe_supersede.h
- ListUsers
- installPage
- pinstate.go
- gui/session.go
- runMonitor
- changelog.go
- FileTarget
- fcntl
- copyRegularAt
- fileedit_symlink_test.go
- time.Time
- btrfsLayoutFrom
- Network
- NewDumpHook
- runWizard
- NetEventType
- process_vm_readv.c
- GuardEvent
- expandPlaceholders
- newTestTarget
- runUninstall
- Monitor
- runNetworkGuard
- Finder
- liveSession
- highlight_test.go
- daemon mode (fscrypt + whitelist lifecycle)
- FORWARD temporary -w/-b whitelist grants
- compileGlobNames
- DaemonEvent
- netinject.c
- NetGuard
- check-compatibility.sh
- BPF verifier 1M-insn complexity budget
- Inode-based (dev:ino) executable identity
- gui.go
- policyEditor
- inode_dev
- Well-known bypass classes
- runLiveEdit
- jit_provenance.c
- inspector_probe.c
- ptrace_race.c
- check_watched
- Build & sign release assets (reusable workflow)
- mutedText
- KernelDev
- mountTable
- NetworkMonitorRepository
- app-listener Terminal Demo Recording
- IntegrationSuite
- IntegrationSuite
- validateConfigText
- IntegrationSuite
- IntegrationSuite
- Regenerate eBPF bindings in pinned builder image
- CopyTreeWithProgress
- filevault_test.go
- WebSocket /ws client
- launch_probe.c
- install.sh
- mc_tag.h
- User
- make test-integration (rootful Docker suite)
- IntegrationSuite
- supersede_probe.c
- NewGuard
- trace-app-libs.sh
- IntegrationSuite
- go_pkg_github_com_testcontainers_testcontainers_go
- Bun TMPDIR redirect with reserved .bun-* name
- Daemon self-protection (selfguards.go)
- github.com/Virgula0/app-listener
- network-monitor mode
- uninstall subcommand
- fakeGuardRepo
- gui/install_test.go
- PlanDaemonUnitUpdate
- pathCache
- WithBottomBar
- net_tester/main.go
- .eventTypeAllowed
- newTestVetter
- .Unlock
- daemon: permanent protection with encryption at rest
- ConfinedEntry
- EventType
- layout.go
- BinaryStat
- IntegrationSuite
- Installation
- os.File
- graphify-refresh.sh
- time.Duration
- activity
- main_test.go
- inode_dev.h
- feed.go
- shell
- writeWithin
- harnessSuite
- downloadFile
- .refreshOnce
- IntegrationSuite
- Limitations
- Troubleshooting
- bpfLsmListed
- mountsUnitTest
- Daemon configuration
- edit-protected
- Vault
- configedit_test.go
- Compatibility
- How it works
- .netIntegrityContainer
- launch_rule_at
- loadDaemonConfig
- golang.org/x/sys/unix.Timespec
- Binary trust
- guard: block file access by program
- network-guard: block network access by program
- .runEditorHarness
- Development
- app-listener documentation
- Self-updating apps and the catalog
- network-monitor: see what a program talks to
- batchInodeRows
- monitor: see what touches a path
- kernel-prof.sh
- Execute
- MulticallError
- .populateWatchInodes
- Desktop GUI

## God Nodes (most connected - your core abstractions)
1. `Guard` - 117 edges
2. `Config` - 96 edges
3. `GuardInodeKey` - 96 edges
4. `Load()` - 90 edges
5. `IntegrationSuite` - 82 edges
6. `absPath()` - 81 edges
7. `writeConfig()` - 69 edges
8. `IntegrationSuite` - 68 edges
9. `shell` - 65 edges
10. `BinaryEntry` - 60 edges

## Surprising Connections (you probably didn't know these)
- `Well-known bypass classes` --references--> `ptrace_race POC`  [INFERRED]
  CLAUDE.md → integrationtests/exploits/README.md
- `lib_probe reserved library POC` --semantically_similar_to--> `Bun TMPDIR redirect with reserved .bun-* name`  [INFERRED] [semantically similar]
  integrationtests/exploits/README.md → README.md
- `swap_benign / swap_reader in-place binary swap POC` --semantically_similar_to--> `Binary ledger (/etc/app-listener/binaries.db)`  [INFERRED] [semantically similar]
  integrationtests/exploits/README.md → README.md
- `runConfigEditor()` --calls--> `EditText()`  [INFERRED]
  cmd/functions/install/steps.go → internal/install/editor.go
- `resolveAssets()` --calls--> `asset()`  [INFERRED]
  cmd/functions/update/update.go → internal/gui/icons.go

## Import Cycles
- None detected.

## Communities (203 total, 32 thin omitted)

### Community 0 - "daemon/daemon.go"
Cohesion: 0.10
Nodes (13): tempJournalRow, launchRule, trustEvent, textAnswer, btrfsIoctlFsInfoArgs, daemonStatsMsg, guardStatsMsg, hint (+5 more)

### Community 1 - "testing.T"
Cohesion: 0.02
Nodes (127): plantBun(), TestOpenBunTmp_KeepsVettedDir(), TestDaemonSelfBaselineEventsNeverWrites(), TestGroupUnlockSelfEventsIsBaselinePlusMknod(), TestValidateFlagsEditConfig(), validateFlags(), TestValidateForwardFlags(), validateForwardFlags() (+119 more)

### Community 2 - "IntegrationSuite"
Cohesion: 0.09
Nodes (10): configBinaryLines(), IntegrationSuite, parseDaemonEvents(), probeExeAs(), procGateMatrix(), procReadProbe(), rawExec(), shQuote() (+2 more)

### Community 3 - "GenerateConf"
Cohesion: 0.05
Nodes (67): logRefreshChange(), TestLogRefreshChange(), TestDiffMergePreservesExistingAndParses(), groupedDiscordConf(), TestGroupedSectionAddressedByEncryptionRoot(), TestPatchCatalogSectionGroupedConfig(), TestSetSectionWhitelistPreservesGroupStructure(), LibraryBlock (+59 more)

### Community 4 - "daemonconfig_test.go"
Cohesion: 0.08
Nodes (73): Load(), ResolvePendingPaths(), hasBinary(), hasPendingBinary(), TestElectronAppsBlockParsed(), TestElectronAppsBlockRejectsMalformedLines(), TestEncryptionGroupsSkipsLibDirs(), TestInspectorsBlockBounded() (+65 more)

### Community 5 - "networkguard.bpf.c"
Cohesion: 0.11
Nodes (34): code_vouched(), emit_gate(), env_loader(), env_step(), exe_inode_of(), fill_current_image(), fill_inode_key(), get_task_exe_key() (+26 more)

### Community 6 - "monitor.bpf.c"
Cohesion: 0.15
Nodes (32): emit_event(), emit_event_kern(), fill_identity(), fill_path(), read_path(), trace_do_faccessat(), trace_do_sendfile(), trace_do_splice() (+24 more)

### Community 7 - "rcUser"
Cohesion: 0.14
Nodes (31): revertBunLaunchers(), rcTarget, bunTargets(), EnsureAddKeysToAgent(), EnsureBunLauncherEnv(), ensureRCBlock(), EnsureSSHAgentEnv(), openOwnedRegular() (+23 more)

### Community 8 - "guardUnitTest"
Cohesion: 0.05
Nodes (24): addErrKind, bpfGuardEvent, classifyAddErr(), decodeGuardEvent(), eventMask(), fsGateLabel(), guardModeKey(), parseGuardEvent() (+16 more)

### Community 9 - "mkdirs"
Cohesion: 0.13
Nodes (36): newWatchPattern(), newFakeInotify(), TestWatchSetBinaryWriteTriggersNotMkdir(), TestWatchSetFloodIsCapped(), TestWatchSetLibDirAppearanceTriggers(), TestWatchSetMissingRootWatchesDeepestAncestor(), TestWatchSetOverflowRescansAndRebuilds(), TestWatchSetSymlinkCreateTriggers() (+28 more)

### Community 10 - "SanitizeText"
Cohesion: 0.09
Nodes (18): multicallAdmissible(), newBinaryVetter(), openBinaryVetter(), resolveInspectors(), TestResolveInspectorsKeepsOnlyAdmitted(), startTrustGuard(), trustStartupError(), runHeadless() (+10 more)

### Community 11 - "buildTrustedSet"
Cohesion: 0.09
Nodes (29): inspectorPaths(), TestInspectorPathsOnlyRootPlaced(), buildTrustedSet(), TestBuildTrustedSetExcludesLibraryClosure(), closureRejections(), inGuardedTree(), keys(), TestWarnUntrustedLibs_OnlyWhatTheKernelRefuses() (+21 more)

### Community 12 - "UI"
Cohesion: 0.08
Nodes (20): installAction, InstallMode, InstallRequest, Options, cancelUI, InstallHooks, Run(), ShowDiff() (+12 more)

### Community 13 - "procstats.go"
Cohesion: 0.15
Nodes (13): Stats, parseKBValue(), parseStat(), parseStatus(), Read(), TestReadConcurrentRace(), TestStatsFieldsNonZero(), TestStatsValuesMonotonic() (+5 more)

### Community 14 - "os.FileMode"
Cohesion: 0.09
Nodes (29): dupFD(), installConfig(), installFile(), installFileAs(), installServices(), installSSHAgent(), installSSHAgentEnv(), installUnit() (+21 more)

### Community 15 - "SelectOrphans"
Cohesion: 0.20
Nodes (17): modifiedContextWithSource(), appProtectorSet(), CleanOrphans(), isAppProtector(), listPolicies(), listProtectors(), orphanedPolicies(), orphanedProtectors() (+9 more)

### Community 16 - "IntegrationSuite"
Cohesion: 0.10
Nodes (9): attrOpCase, eventsForPath(), IntegrationSuite, guardDeltaEvents(), guardEventTypes(), le64HexKey(), parseGuardEvents(), guardEvent (+1 more)

### Community 17 - "guard_trust.bpf.c"
Cohesion: 0.11
Nodes (51): caller_updates(), current_exe_flags(), current_exe_inode(), current_exe_key(), current_writer_bits(), dentry_is_dir(), dentry_parent(), deny_if_protected() (+43 more)

### Community 19 - "fileEditModel"
Cohesion: 0.08
Nodes (12): applyNewFileMeta(), clipLabel(), humanSize(), listEntries(), clamp(), sortNodes(), chownUser, fileEditMode (+4 more)

### Community 20 - "networkmonitor.bpf.c"
Cohesion: 0.31
Nodes (16): BPF_KRETPROBE(), emit_event(), get_netns(), get_socket_proto(), is_watched_binary(), read_inet_addr(), trace_accept(), trace_accept4() (+8 more)

### Community 21 - "usecase/daemon_test.go"
Cohesion: 0.12
Nodes (56): NewDaemonUseCase(), partitionEncryptionRoots(), newFakeVault(), resource(), startDaemon(), TestDaemonUseCaseConstructorMismatch(), TestDaemonUseCaseEventsMerged(), TestDaemonUseCaseGrantEditAccess() (+48 more)

### Community 22 - "Resource"
Cohesion: 0.06
Nodes (12): ownerPaths(), fakeDaemon, Resource, GuardRepository, Vault, concurrencyLimit(), encryptionRootOf(), DaemonUseCase (+4 more)

### Community 23 - "auth.go"
Cohesion: 0.12
Nodes (31): Hash(), HashFileExists(), LoadHashFile(), OriginOf(), parseHash(), RemoveHashFile(), hashFileInode(), TestHashFileLifecycle() (+23 more)

### Community 24 - "NetGuardEvent"
Cohesion: 0.10
Nodes (10): snapNetGuard, NetGuardEvent, NetworkGuardRepository, protoString(), NetworkGuardUseCase, NewNetworkGuardUseCase(), TestNetworkGuardUseCaseStartError(), netGuardEventLine (+2 more)

### Community 25 - "guard.bpf.c"
Cohesion: 0.12
Nodes (53): add_inode_to_guard(), chmod_only_drops_write(), count_degrade(), discover_guarded_parent(), emit_process_denial(), emit_with(), event_is_read_class(), evict_inode_from_guard() (+45 more)

### Community 26 - "Ledger"
Cohesion: 0.07
Nodes (25): newHashCache(), ds(), TestHashCacheFlushKeepsOnlyThisRunsRows(), candidates(), confirmable(), confirmAll(), describe(), hashNow() (+17 more)

### Community 27 - "networkPage"
Cohesion: 0.17
Nodes (8): binaryList, networkPage, watchTab, newBinaryList(), entryPaths(), newNetGuardTab(), optionHint(), newWatchTab()

### Community 28 - "inodeChain"
Cohesion: 0.29
Nodes (12): inodeChain(), isProcFD(), mountinfoShowsSbRoot(), mountShowsSbRoot(), keyOf(), mkTree(), TestInodeChainDirIncludesItself(), TestInodeChainFollow() (+4 more)

### Community 29 - "github.com/charmbracelet/bubbles/textarea.Model"
Cohesion: 0.08
Nodes (45): class, lineStyles, viewportState, classOf(), column(), dedent(), Edit(), leadingIndent() (+37 more)

### Community 30 - "journal"
Cohesion: 0.15
Nodes (6): fakeGuard, fakeMonitor, fakeNet, fakeNetGuard, fakeNetMonitor, journal

### Community 31 - "serve.go"
Cohesion: 0.07
Nodes (35): basicAuth(), configureWebSocket(), decodeEndsAtEOF(), isInteractiveTerminal(), NewEventFanout(), newServeListener(), newServeServer(), parseResize() (+27 more)

### Community 32 - "bpfstats/main.go"
Cohesion: 0.09
Nodes (31): CheckBPFLSM(), confirmMasterKeyOverwrite(), runBPFCheck(), runGenKey(), runOneShotMode(), ensureMasterKey(), checkKeyLen(), GenerateMasterKey() (+23 more)

### Community 33 - "update.go"
Cohesion: 0.17
Nodes (23): applyUpdate(), assetsFor(), checkUpdatePreconditions(), confirmUnit(), confirmUpdate(), downloadAndVerify(), FetchLatest(), FetchVerified() (+15 more)

### Community 34 - "filevault.go"
Cohesion: 0.18
Nodes (17): classifyRegularFileTarget(), clearRecovery(), deriveFileVaultSubkey(), EnsureRecoverySidecarPlaceholder(), fileVaultRecoverPath(), openFileVaultWithMasterKey(), recoverFileInPlace(), rejectSidecarSymlink() (+9 more)

### Community 35 - "github.com/spf13/cobra.Command"
Cohesion: 0.08
Nodes (31): AddServeFlags(), credentialFlagState(), isLoopbackHost(), ParseServeFlags(), requestedBoolFlag(), validateEnabledServe(), validateServeAddress(), TestParseServeFlagsRequiresCredentials() (+23 more)

### Community 37 - "BinaryEntry"
Cohesion: 0.07
Nodes (28): binaryVerifyState, deferredBinary, HashStore, sharedHashEntry, BinaryEntry, TestCheckBinaryEventsReadOnlyAcceptsUnrestricted(), TestCheckBinaryEventsReadOnlyRejectsRestriction(), checkBtrfsKeys() (+20 more)

### Community 38 - "controlServer"
Cohesion: 0.12
Nodes (18): readConfigPut(), TestReadGrantRequestConfig(), readAuthRequest(), readClientLines(), readForwardRequest(), readGrantRequest(), readPathLines(), streamEvents() (+10 more)

### Community 39 - "fileedit_model_test.go"
Cohesion: 0.13
Nodes (32): OpenRegularNoFollow(), defaultNewFileMode(), fileModeToUnixMode(), isBinaryContent(), isOctalMode(), newFileEditModel(), inodeOf(), maxLineWidth() (+24 more)

### Community 40 - "Guard"
Cohesion: 0.07
Nodes (9): exeRow, rootHandle, tempMaskOp, tempRow, openRoot(), SharedPinDegraded(), TestNoteRootMountFollowsMountID(), Guard (+1 more)

### Community 41 - "IntegrationSuite"
Cohesion: 0.20
Nodes (3): exploitTest, newEventTypes(), IntegrationSuite

### Community 42 - "absPath"
Cohesion: 0.06
Nodes (5): IntegrationSuite, IntegrationSuite, inNosuidNamespace(), IntegrationSuite, absPath()

### Community 44 - "IntegrationSuite"
Cohesion: 0.16
Nodes (7): netMonitorEvent, IntegrationSuite, fmtEvents(), IntegrationSuite, netMonitorDeltaEvents(), netMonitorHasTypes(), parseNetMonitorEvents()

### Community 45 - "update_test.go"
Cohesion: 0.09
Nodes (18): fetchReleases(), parseChecksum(), parsePublicKey(), rsaPublicKeyPEM(), signedChecksum(), TestFetchReleases(), TestFetchReleasesHTTPError(), TestFilterChannel() (+10 more)

### Community 46 - "daemon_multicall_test.go"
Cohesion: 0.17
Nodes (10): decodeMulticall(), IntegrationSuite, multicallBypasses(), multicallConfig(), uutilsStrays(), vettedGroupIndex(), tailLast(), mcBypass (+2 more)

### Community 47 - "stdio"
Cohesion: 0.07
Nodes (4): put(), wait_then_read(), hold_forever(), main()

### Community 48 - "golang.org/x/sys/unix.Stat_t"
Cohesion: 0.17
Nodes (12): placedWalk, binaryStatOf(), confineBelow(), inUserHome(), GlobSystemPlaced(), placedDir(), placedMatches(), readlinkFd() (+4 more)

### Community 49 - "sync.Mutex"
Cohesion: 0.07
Nodes (20): drain(), newCatalogRefresher(), readEvents(), swapFile(), TestDebouncerCoalescesAndBoundsWait(), TestDebouncerMinGap(), newControlManager(), startControlManager() (+12 more)

### Community 50 - "LibraryClosure"
Cohesion: 0.10
Nodes (19): ldConfParser, defaultLibDirs(), elfInterp(), fileExists(), isScript(), ldConfInclude(), ldSoConfDirs(), LdSoPreloadPaths() (+11 more)

### Community 51 - "runOfflineEdit"
Cohesion: 0.14
Nodes (14): TestPickEncryptedDirSingleSkipsPicker(), pickEncryptedDir(), runOfflineEdit(), askDecrypt(), decryptStep(), pickDirsToDecrypt(), TestExistingEncryptedDirsNilPaths(), TestExistingEncryptedDirsSkipsPlain() (+6 more)

### Community 52 - ".FilterExistingWhitelist"
Cohesion: 0.10
Nodes (20): confSafeMatch(), homeMatchConfined(), symlinkStaysInParent(), TestExpandLibDirWritersGlob(), TestFilterExistingWhitelist(), TestFilterExistingWhitelistGlob(), TestFilterExistingWhitelistMixed(), TestFilterExistingWhitelistSkipsOpaqueMulticall() (+12 more)

### Community 53 - "Config"
Cohesion: 0.06
Nodes (65): systemPatterns(), systemRule(), appendSectionsAndEdit(), collectDiffAdditions(), configPaths(), TestUncoveredCandidates(), uncoveredCandidates(), grantNewBinaries() (+57 more)

### Community 54 - "github.com/charmbracelet/bubbletea.Cmd"
Cohesion: 0.06
Nodes (25): diffModel, editorModel, formatType(), listenForEvents(), syncViewport(), tickStats(), formatAddr(), formatNetProto() (+17 more)

### Community 55 - "TrustGuard"
Cohesion: 0.05
Nodes (24): setConfinedHomes(), tempGrant, trustHook, trustRows, deleteInoKeys(), deleteResKeys(), setSupersedeTrust(), liveTrust() (+16 more)

### Community 56 - ".showSized"
Cohesion: 0.12
Nodes (10): installUI, multiAnswer, resizeWatch, sizedDialog, scrollable(), selectable(), ask(), diffText() (+2 more)

### Community 57 - "systemd.go"
Cohesion: 0.08
Nodes (40): applyDiffAdditions(), restoreDaemonAfterDiffAbort(), applyLiveRefresh(), buildBinaryIfNeeded(), installBinaryOnly(), runUpdateCatalogOnly(), softenAutomatedRefreshErr(), TestApplyLiveRefreshDeliversReload() (+32 more)

### Community 58 - "GuardInodeKey"
Cohesion: 0.03
Nodes (51): engine, refusedReplacements, resInfo, taintOwner, TestNearestRootClaims(), TestTieRank(), nextLabel(), TestDeliverLabelsCrossResourceProcessGates() (+43 more)

### Community 60 - "buildOneGuard"
Cohesion: 0.12
Nodes (26): buildOneGuard(), eventFilterOptions(), TestEventFilterOptions(), newSelfGuards(), selfProtectSpecs(), startSelfGuard(), startSelfGuards(), pinCfg (+18 more)

### Community 61 - "icon"
Cohesion: 0.04
Nodes (13): appTheme, iconTileRenderer, navItem, navItemRenderer, panelRenderer, pillRenderer, asset(), coloredIcon() (+5 more)

### Community 62 - "exe_supersede.h"
Cohesion: 0.17
Nodes (18): exe_last_link_key(), exe_note_freed(), exe_note_supersede(), exe_seq_now(), exe_stamp_exec(), exe_stamp_fork(), exe_stamp_free(), exe_stamp_lost() (+10 more)

### Community 63 - "ListUsers"
Cohesion: 0.09
Nodes (24): pathCovered(), TestPathCovered(), cleanOrphanedFscrypt(), isInsidePath(), TestAskSSHAgentUsersNoGuardedSSH(), askSSHAgentUsers(), cleanOrphanedMetadata(), detectSSHAgentUnits() (+16 more)

### Community 64 - "installPage"
Cohesion: 0.10
Nodes (9): binSource, installedBinary, installPage, runView, stageItem, stageState, hashReader(), readInstalled() (+1 more)

### Community 65 - "pinstate.go"
Cohesion: 0.17
Nodes (23): ensurePinStateFilePlaceholder(), ensurePlaceholder(), pinOwnerLikelyAlive(), readPinState(), recoverPinState(), pinStateInode(), TestPinOwnerLikelyAlive(), TestPinOwnerLikelyAliveDeadPID() (+15 more)

### Community 66 - "gui/session.go"
Cohesion: 0.07
Nodes (40): ancestorEnv(), invokingHome(), invokingUID(), ownedBy(), readOwnedFile(), restoreSessionEnv(), statusField(), TestDisplayValue() (+32 more)

### Community 67 - "runMonitor"
Cohesion: 0.09
Nodes (21): BuildEBPFTargets(), CheckEBPF(), RawTarget, isSubDir(), MakeDisplayPaths(), ParseEventsFlag(), pathWithinDir(), ResolveTargets() (+13 more)

### Community 68 - "changelog.go"
Cohesion: 0.17
Nodes (23): changelogText(), cumulativeChangelog(), newChangelogModel(), newReleaseWindow(), releasesSince(), showChangelog(), rel(), tags() (+15 more)

### Community 69 - "FileTarget"
Cohesion: 0.09
Nodes (15): filesPage, specFromForm(), specSummary(), EventTypes(), newExeResolver(), Files, FileSpec, FileTarget (+7 more)

### Community 70 - "fcntl"
Cohesion: 0.09
Nodes (9): main(), report(), denied(), main(), dump_via_debugfs(), main(), probe_open(), denied() (+1 more)

### Community 71 - "copyRegularAt"
Cohesion: 0.31
Nodes (10): copyDirContents(), copyEntry(), copyRegularAt(), copySymlinkAt(), copyTreeRoot(), fchownFromInfo(), futimesFromInfo(), readDirNames() (+2 more)

### Community 72 - "fileedit_symlink_test.go"
Cohesion: 0.18
Nodes (15): openRootT(), TestWriteFileKeepMeta(), TestWriteFileKeepMetaRefusesSymlink(), mustRead(), mustWrite(), swappedParent(), TestChmodRefusesSymlinkedParent(), TestCreateEntryRefusesPlantedSymlink() (+7 more)

### Community 73 - "time.Time"
Cohesion: 0.10
Nodes (20): refreshAdmits(), TestRefreshAdmitsHomeMatchesAndRootPlacedNamesOnly(), addSystemBinary(), catalogPatterns(), catalogWatchPlan(), inHome(), isSymlink(), newWatchSet() (+12 more)

### Community 74 - "btrfsLayoutFrom"
Cohesion: 0.14
Nodes (12): BtrfsLayout, fakeTypes, typeSource, ensureBtrfsLayout(), btrfsLayoutFrom(), findMember(), memberOffset(), ResolveBtrfsLayout() (+4 more)

### Community 75 - "Network"
Cohesion: 0.16
Nodes (5): netGuardPill(), clonePolicy(), NetExtra, NetPolicy, Network

### Community 76 - "NewDumpHook"
Cohesion: 0.11
Nodes (17): applyVerbosity(), logHook, logLine, VerboseLevel, BridgeStdLog(), formatEntry(), NewDumpHook(), sortedKeys() (+9 more)

### Community 77 - "runWizard"
Cohesion: 0.11
Nodes (21): askBunTmpdirUsers(), gatherPerUserSetup(), setupBunTmpdir(), promptEditPassword(), checkInstalledMatchesSource(), checkRunningBinaryMatchesInstalled(), cleanupBackups(), deploy() (+13 more)

### Community 78 - "NetEventType"
Cohesion: 0.07
Nodes (18): snapNetMonitor, NetBpfEvent, NetEvent, NetEventType, newNetEventChecks(), Cstr(), FormatAddr(), NetEventTypes() (+10 more)

### Community 79 - "process_vm_readv.c"
Cohesion: 0.18
Nodes (6): dump_memory(), find_all_rw_regions(), find_heap(), find_victim_by_fd(), main(), spawn_less_pty()

### Community 80 - "GuardEvent"
Cohesion: 0.06
Nodes (22): newTestSession(), TestRunSessionIdleTimeoutResetsOnPing(), TestRunSessionStreamsEventsAsActivity(), holdHeadless(), HeadlessLine(), runGuardHeadless(), snapGuard, TestDispatchFullChannelNeverDropsDenials() (+14 more)

### Community 81 - "expandPlaceholders"
Cohesion: 0.14
Nodes (8): symlinkTargetGlob(), expandPlaceholders(), TestExpandWhitelist(), TestExtraWatchPathsForSkipsMissing(), LibDirGlob, TrustGlob, reservedGlobs(), splitTrustGlob()

### Community 82 - "newTestTarget"
Cohesion: 0.21
Nodes (17): Target, FileEngines, NewFiles(), newTestTarget(), since(), TestAddRefusesOverlap(), TestApplyFromBlacklistFailureKeepsOldGuard(), TestApplyFromBlacklistOverlaps() (+9 more)

### Community 83 - "runUninstall"
Cohesion: 0.13
Nodes (25): runDiffCatalog(), parseMaintenanceFlags(), runInstall(), runMaintenanceMode(), TestValidateMaintenanceFlags(), validateMaintenanceFlags(), deletePostBackups(), restoreBackups() (+17 more)

### Community 84 - "Monitor"
Cohesion: 0.15
Nodes (3): NewMonitor(), Monitor, watchedInode

### Community 85 - "runNetworkGuard"
Cohesion: 0.11
Nodes (22): ParseNetEventsFlag(), computeGuardBinaries(), confirmUnsafe(), discoverInfra(), resolveGuardMode(), runNetworkGuard(), runTUI(), WithSuspectsFrom() (+14 more)

### Community 86 - "Finder"
Cohesion: 0.14
Nodes (19): findButton, span, findAll(), Finder, lowerRunes(), NewFinder(), newFindInput(), overlay() (+11 more)

### Community 87 - "liveSession"
Cohesion: 0.18
Nodes (12): parseErr(), TestLiveSessionConfigAndPut(), TestLiveSessionWatchAndPinger(), chooseResources(), confirmUserPlaced(), displayEntries(), eventsLabel(), forwardBinaries() (+4 more)

### Community 88 - "highlight_test.go"
Cohesion: 0.16
Nodes (12): Highlighter, lexerFor(), NewHighlighter(), shebangLexer(), sampleGo(), TestDaemonConfTokens(), TestHighlighterOnlyChangesColors(), TestHighlighterUnstyledMatchesTextareaView() (+4 more)

### Community 89 - "daemon mode (fscrypt + whitelist lifecycle)"
Cohesion: 0.18
Nodes (10): Daemon catalog refresh (catalogrefresh.go / catalogwatch.go), exe_supersede.h superseded-key tracking, glob_plant reserved glob name POC, launch_probe risky-launch POC, preload_wait code-suspect POC, supersede_probe reused-inode POC, Daemon boot ordering (Type=notify, Before=user sessions), daemon mode (fscrypt + whitelist lifecycle) (+2 more)

### Community 90 - "FORWARD temporary -w/-b whitelist grants"
Cohesion: 0.21
Nodes (9): CONFIG post-AUTH daemon.conf editing (configedit.go), FORWARD temporary -w/-b whitelist grants, Daemon self-key (GUARD_ALLOW_ROOT) bypass, swap_benign / swap_reader in-place binary swap POC, edit-protected --edit-config, edit-protected subcommand, edit-protected live mode, edit-protected --forward (+1 more)

### Community 91 - "compileGlobNames"
Cohesion: 0.17
Nodes (9): GlobChild, statFunc, compileGlobNames(), globText(), resolveBits(), resolveChildren(), statKey(), TestCompileGlobNames_RepeatedPatternOrsBits() (+1 more)

### Community 92 - "DaemonEvent"
Cohesion: 0.08
Nodes (37): ServeConfig, UIDResolver, NewUIDResolver(), drainEphemeral(), newDaemonModel(), runDaemonUI(), runHeadless(), runServedTUI() (+29 more)

### Community 93 - "netinject.c"
Cohesion: 0.17
Nodes (10): attach_child(), await_exec(), do_connect(), errname(), main(), sleep_ms(), traceme(), do_connect() (+2 more)

### Community 95 - "NetGuard"
Cohesion: 0.09
Nodes (11): foreignPinnedLink(), eventTypeKey(), mcNameKey(), statInodeKey(), judgePasses(), parseProcs(), TestJudgePasses(), TestParseProcs() (+3 more)

### Community 96 - "check-compatibility.sh"
Cohesion: 0.38
Nodes (11): fail(), load_config(), lsm_bpf_instructions(), note(), pass(), read_config(), require_config(), section() (+3 more)

### Community 97 - "BPF verifier 1M-insn complexity budget"
Cohesion: 0.20
Nodes (8): PR vs base verifier cost comparison, Verifier gate workflow, bpf_loop() helper, tools/bpfstats verifier cost tool, static __noinline BPF-to-BPF helpers with args struct, rename-over-guarded-file bypass, xres_move cross-resource move POC, Kernel floor (5.8 / 5.10 / 5.17)

### Community 98 - "Inode-based (dev:ino) executable identity"
Cohesion: 0.18
Nodes (7): CheckBPFLSM preflight, internal/infrastructure shared eBPF kernel, Ports & adapters layering (cmd -> usecase -> repository -> engines), jit_provenance ledger (guard_jit_origin), daemon --check / --verifier-only preflight, guard mode (LSM whitelist/blacklist), network-guard mode

### Community 99 - "gui.go"
Cohesion: 0.17
Nodes (12): needsBinary(), RunWith(), TestNeedsBinary(), TestRunWithRefusesDeployWithoutBinary(), TestRunWithRefusesNonHeadlessBinary(), sanityCheckBinary(), TestSanityCheckBinaryRefusesOtherMachine(), Mode (+4 more)

### Community 100 - "policyEditor"
Cohesion: 0.10
Nodes (14): netGuardTab, policyEditor, policyHooks, policyTab, guardPill(), newPolicyTab(), verbFor(), eventsOrAll() (+6 more)

### Community 101 - "inode_dev"
Cohesion: 0.18
Nodes (15): exe_is_superseded(), exe_refused(), inode_dev(), check_and_emit_args(), code_suspect(), current_is_inspector(), exe_tag_bits(), fill_path() (+7 more)

### Community 102 - "Well-known bypass classes"
Cohesion: 0.24
Nodes (9): btrfs copy-ioctl gate (file_ioctl / file_ioctl_compat), Well-known bypass classes, btrfs_search POC, io_uring POC, open_by_handle_at POC, process_vm_readv POC (remaining monitor gap), raw_block_device POC, statonly stat/statx metadata leak (+1 more)

### Community 103 - "runLiveEdit"
Cohesion: 0.09
Nodes (27): auditAfterEdit(), auditAfterEditWithConfig(), auditEntry(), auditGroupCoverage(), auditTree(), coveredByWatch(), findResource(), groupWatchPaths() (+19 more)

### Community 104 - "jit_provenance.c"
Cohesion: 0.47
Nodes (8): copy_into(), main(), make_temp_copy(), mode_hold(), mode_memfd(), mode_passfd(), mode_self(), report_dlopen()

### Community 105 - "inspector_probe.c"
Cohesion: 0.56
Nodes (8): check_attach(), check_exe(), check_mem(), check_root(), check_vmread(), main(), report(), traced_exec()

### Community 106 - "ptrace_race.c"
Cohesion: 0.31
Nodes (5): main(), peek_buffer(), run_tracer(), run_victim(), wait_for_file()

### Community 107 - "check_watched"
Cohesion: 0.24
Nodes (14): admit_refusal(), check_watched(), emit_event(), guard_net_socket_bind(), guard_net_socket_connect(), guard_net_socket_listen(), guard_net_socket_recvmsg(), guard_net_socket_sendmsg() (+6 more)

### Community 108 - "Build & sign release assets (reusable workflow)"
Cohesion: 0.31
Nodes (8): Build & sign release assets (reusable workflow), app-listener-gui best-effort build, Promote to stable release workflow, Release workflow (pre-release per push to main), Release promote flow, make check-compatibility, scripts/install.sh one-line installer, update self-update subcommand

### Community 109 - "mutedText"
Cohesion: 0.07
Nodes (28): binaryRow, clickPanel, counterStrip, homePage, iconTile, panel, pill, statusTile (+20 more)

### Community 110 - "KernelDev"
Cohesion: 0.10
Nodes (24): fdStat, physicalParent(), statFD(), isFuseType(), parseMajorMinor(), TestParseMajorMinorMatchesStatEncoding(), TestVouchedDevs(), TestVouchedDevsFailsClosed() (+16 more)

### Community 111 - "mountTable"
Cohesion: 0.16
Nodes (9): MountInfo, mountinfoUnitTest, MountByID(), parseMountDev(), ParseMountinfo(), PathUnder(), ReadMountinfo(), unescapeMountinfo() (+1 more)

### Community 112 - "NetworkMonitorRepository"
Cohesion: 0.12
Nodes (13): NetworkMonitorRepository, NewGuardUseCase(), NetworkMonitorUseCase, NewNetworkMonitorUseCase(), newFakeMonitorRepo(), newFakeNetworkGuardRepo(), newFakeNetworkMonitorRepo(), TestGuardUseCaseLifecycle() (+5 more)

### Community 113 - "app-listener Terminal Demo Recording"
Cohesion: 0.39
Nodes (8): BPF Build Pipeline (bpf2go + go build), CLI Subcommand Help (guard, install, monitor, network-guard, network-monitor), app-listener Terminal Demo Recording, Guard TUI Blacklist Mode (eBPF LSM), Guard TUI Whitelist Mode (eBPF LSM), File System Monitor TUI (eBPF), Startup Banner and eBPF Availability Check, TUI Resource Footer (Mem / CPU / Evt/s)

### Community 116 - "validateConfigText"
Cohesion: 0.13
Nodes (14): collectFilesystemPrereqs(), askEncryption(), discordCatalogEntry(), TestAskEncryptionSkipsNeedEncryptionFalse(), TestCollectFilesystemPrereqsNoPanic(), TestCollectFilesystemPrereqsSkipsRegularFiles(), TestGeneratedLibraryBlockRoundTrips(), TestSectionsFromCandidatesSkipsMissingWatchSubPaths() (+6 more)

### Community 119 - "Regenerate eBPF bindings in pinned builder image"
Cohesion: 0.33
Nodes (5): CI workflow (lint + test on non-draft PRs), CI lint job (make lint), CI test job (make test), golangci-lint v2 configuration, Issue #45 (clang 14 guard_path_rename over budget)

### Community 120 - "CopyTreeWithProgress"
Cohesion: 0.20
Nodes (5): migrationTarget, classifyMigrationTarget(), CopyTreeWithProgress(), TestCopyTreeProgress(), treeSize()

### Community 121 - "filevault_test.go"
Cohesion: 0.18
Nodes (26): looksLikeFileVaultRecord(), newFileVaultAEAD(), openFileVault(), sealFileVault(), inode(), mustSubkey(), TestDeriveFileVaultSubkeyDeterministicAndSeparated(), TestIsEncryptedDispatchesToFileVaultForRegularFiles() (+18 more)

### Community 122 - "WebSocket /ws client"
Cohesion: 0.29
Nodes (3): AppListener Shared Session page, WebSocket /ws client, --serve WebSocket mirroring

### Community 124 - "launch_probe.c"
Cohesion: 0.53
Nodes (4): copy(), is_verb(), main(), usage()

### Community 125 - "install.sh"
Cohesion: 0.57
Nodes (6): die(), download(), info(), install.sh script, usage(), warn()

### Community 126 - "mc_tag.h"
Cohesion: 0.24
Nodes (10): mc_attest(), mc_basename(), mc_leader_start(), mc_stamp_fork(), mc_stamp_free(), mc_tag_bits(), netg_exec_applet(), netg_task_free() (+2 more)

### Community 127 - "User"
Cohesion: 0.12
Nodes (27): entryWriters(), reserveChain(), TestGlobBuilder_LibBitsPerEntry(), underResource(), bunEntryConfigured(), bunLaunchersFor(), globBuilder, ParseGlobName() (+19 more)

### Community 128 - "make test-integration (rootful Docker suite)"
Cohesion: 0.40
Nodes (4): make test-integration (rootful Docker suite), app-listener docker-compose service, ptrace_race POC, Docker usage (docker compose)

### Community 130 - "supersede_probe.c"
Cohesion: 0.90
Nodes (4): answer(), main(), read_secret(), run_child()

### Community 131 - "NewGuard"
Cohesion: 0.19
Nodes (11): resolveGuardConfig(), runGuard(), runGuardTUI(), canonicalPaths(), Mode, modeString(), NewGuard(), WithEagerPopulate() (+3 more)

### Community 132 - "trace-app-libs.sh"
Cohesion: 0.60
Nodes (3): conf_directive_values(), list_conf_binaries(), trace-app-libs.sh script

### Community 134 - "go_pkg_github_com_testcontainers_testcontainers_go"
Cohesion: 0.20
Nodes (3): guardNetTypesForComm(), netGuardBlockedEventCount(), netGuardHasBlockedEvent()

### Community 135 - "Bun TMPDIR redirect with reserved .bun-* name"
Cohesion: 0.50
Nodes (4): lib_probe reserved library POC, Bun TMPDIR redirect with reserved .bun-* name, install --diff-catalog, install TUI wizard

### Community 142 - "fakeGuardRepo"
Cohesion: 0.05
Nodes (19): SupersededBinary, TempBinary, TemporaryGrant, TempRule, planTempBlock(), applyTempGrants(), closeTempBinaries(), TemporaryAccess (+11 more)

### Community 143 - "gui/install_test.go"
Cohesion: 0.14
Nodes (11): stagedBinary, stageLocal(), TestStageLocal(), TestStageLocalRefusedIsRemoved(), entryOf(), BtrfsMounted(), mountinfoDev(), mountinfoHasFstype() (+3 more)

### Community 144 - "PlanDaemonUnitUpdate"
Cohesion: 0.21
Nodes (12): TestDaemonUnitWithMetadataOutput(), DaemonUnitWithMetadataOutput(), ConfirmUnitFunc, InstalledDaemonUnit(), PlanDaemonUnitUpdate(), ShippedDaemonUnit(), fakeBinary(), TestPlanDaemonUnitUpdate() (+4 more)

### Community 145 - "pathCache"
Cohesion: 0.39
Nodes (4): evalSymlinksOrEmpty(), newPathCache(), pathCache, pathCacheEntry

### Community 146 - "WithBottomBar"
Cohesion: 0.21
Nodes (7): decryptDirectories(), TestDecryptDirectoriesEmpty(), TestProgressReader(), bottomBar, newBottomBar(), renderBar(), WithBottomBar()

### Community 147 - "net_tester/main.go"
Cohesion: 0.21
Nodes (7): connFD(), tcpClient(), tcpServer(), tcpServerDelayed(), udpRecvLoop(), udpSendLoop(), waitForMarker()

### Community 148 - ".eventTypeAllowed"
Cohesion: 0.29
Nodes (3): TestEventTypeAllowedConcurrent(), TestSetEventTypesBeforeStartNoRace(), TestSetEventTypesRace()

### Community 149 - "newTestVetter"
Cohesion: 0.44
Nodes (10): newTestVetter(), testBinary(), TestVetterBootstrapRecordsThenChecks(), TestVetterLinkKeyedOnLink(), TestVetterRecordLiveUnderLinkLines(), TestVetterRefusesNewLineAfterBootstrap(), TestVetterResolverOnlyApproved(), TestVetterSeesInPlaceRewriteThroughSharedHashes() (+2 more)

### Community 150 - ".Unlock"
Cohesion: 0.12
Nodes (18): deprovisionKind, isFileVaultCiphertext(), isRegularFileTarget(), classifyDeprovision(), classifyDeprovisionErr(), hasEncryptionPolicy(), isLockedRegularFileErr(), isNotEncryptedErr() (+10 more)

### Community 152 - "daemon: permanent protection with encryption at rest"
Cohesion: 0.17
Nodes (12): Boot ordering, daemon: permanent protection with encryption at rest, Filesystem-wide gates, Flags, Lifecycle, Process inspectors and taint, Quiet metadata denials, Related (+4 more)

### Community 153 - "ConfinedEntry"
Cohesion: 0.29
Nodes (19): ConfinedEntry(), HashStats(), hashedReused(), resetSharedHashes(), TestConcurrentCallersHashOnce(), TestConfinedEntryHashesOncePerFingerprint(), TestNewVerifyStateSeededOnlyByMatchingHash(), TestPersistentStoreFilledAfterSyncAndHit() (+11 more)

### Community 154 - "EventType"
Cohesion: 0.12
Nodes (10): snapMonitor, BpfEvent, EventType, eventUnitTest, FileEvent, parseEvents(), TestAddBinaryActionsRefusesBeforeWriting(), DecodeBpfEvent() (+2 more)

### Community 155 - "layout.go"
Cohesion: 0.45
Nodes (5): adaptiveGrid, gapVBox, cellMin(), rowHeights(), visible()

### Community 156 - "BinaryStat"
Cohesion: 0.29
Nodes (10): BinaryStat, binaryStatOf(), StatBinary(), StatBinaryFile(), durableFS(), fsidOf(), hostVouches(), StatDurable() (+2 more)

### Community 158 - "Installation"
Cohesion: 0.18
Nodes (11): Build from source, Bun-based apps (opencode), Docker, fscrypt prerequisites, Install flags, Installation, One-line installer, SSH agent (+3 more)

### Community 159 - "os.File"
Cohesion: 0.04
Nodes (70): checkFreshDir(), fileID(), openBunTmp(), openDirNoFollow(), renewDir(), TestOpenBunTmp_ReplacesUnvettedDir(), TestOpenBunTmp_SymlinkReplacedTargetUntouched(), bunTmp (+62 more)

### Community 161 - "time.Duration"
Cohesion: 0.13
Nodes (13): procStats, statsSample, cpuTime(), humanBytes(), rssBytes(), New(), recv(), TestSubDeliversMissedTickLate() (+5 more)

### Community 162 - "activity"
Cohesion: 0.25
Nodes (6): newActivity(), TestRingKeepsNewest(), Counts, newRing(), activity, ring

### Community 163 - "main_test.go"
Cohesion: 0.20
Nodes (5): infraContainerPaths(), TestIntegrationSuite(), TestMain(), helperChild(), TestMain()

### Community 164 - "inode_dev.h"
Cohesion: 0.24
Nodes (8): ctx_ptr(), sb_dev(), btrfs_copy_gate(), guard_file_ioctl(), guard_file_ioctl_compat(), guard_inode_free(), is_btrfs_copy_ioctl(), trust_inode_free()

### Community 165 - "feed.go"
Cohesion: 0.08
Nodes (36): actionBar, cell, column, dataTable, eventFeed, feedAction, feedSpec, processView (+28 more)

### Community 166 - "shell"
Cohesion: 0.06
Nodes (33): pageDef, pageID, shell, snapEngines, snapNet, sysStatus, newLogo(), newShell() (+25 more)

### Community 167 - "writeWithin"
Cohesion: 0.22
Nodes (9): atomicWriteAt(), descendNoFollow(), refuseExistingSymlink(), TestWriteWithin_PlainFile(), TestWriteWithin_RefusesSymlinkDirComponent(), TestWriteWithin_RefusesSymlinkTarget(), TestWriteWithin_RejectsEscape(), writeAndSync() (+1 more)

### Community 169 - "downloadFile"
Cohesion: 0.32
Nodes (6): downloadFile(), downloadReleaseFiles(), fetchRelease(), TestDownloadFileHTTPError(), TestDownloadFileMode0700(), TestFetchReleaseByTag()

### Community 170 - ".refreshOnce"
Cohesion: 0.43
Nodes (6): sameConfig(), discordRefreshFixture(), TestRefreshOnceAppliesAndReloads(), TestRefreshOnceKeepsOperatorLinesAndSkipsNoOps(), TestRefreshOnceRestoresTheFileWhenTheReloadFails(), TestRefreshOnceSkipsAConfigChangedSinceLoad()

### Community 172 - "Limitations"
Cohesion: 0.22
Nodes (9): Already-running processes, Binaries you can write, Event masks are not confinement, Interpreters and JITs, Limitations, Project status, Root, Side effects to expect with the daemon (+1 more)

### Community 173 - "Troubleshooting"
Cohesion: 0.22
Nodes (9): A program I use is denied, After a package update replaced a binary, An app breaks and nothing is logged, Backups, Check the guard by hand, Encryption key errors, Is the daemon running and enforcing?, Snapshots or disk backups stopped working (+1 more)

### Community 174 - "bpfLsmListed"
Cohesion: 0.33
Nodes (7): bpfLsmListed(), CheckBPFLSM(), CheckBPFLSMAt(), TestBpfLsmListed(), TestBpfLsmListedMissingFile(), TestCheckBPFLSM(), writeLSM()

### Community 176 - "Daemon configuration"
Cohesion: 0.25
Nodes (8): Daemon configuration, Daemon settings: `[settings]`, Generic Electron apps: `[electron_apps]`, Library trust: `[libraries "<name>"]`, Per-binary event masks, Reloading, Several guarded trees, one vault, Watch sections

### Community 177 - "edit-protected"
Cohesion: 0.25
Nodes (8): edit-protected, Editing the daemon config, Editor keys, Exit audit, Non-interactive writes, Offline and live mode, Temporary access for other programs: `--forward`, The password and the control socket

### Community 178 - "Vault"
Cohesion: 0.07
Nodes (43): backingDeviceUnion(), buildConcurrency(), buildGuards(), configureDaemonLogging(), lockOneRoot(), lockRootRecovering(), newPinGeneration(), prepareDaemonStart() (+35 more)

### Community 179 - "configedit_test.go"
Cohesion: 0.47
Nodes (9): newTestConfigEditor(), readFile(), TestConfigEditorApply(), TestConfigEditorRefusals(), TestConfigEditorRefusesNewMulticallLine(), TestConfigEditorRestoresOnFailedReload(), TestRunConfigEditOverTheSocket(), TestRunConfigEditRefusedDuringALiveSession() (+1 more)

### Community 180 - "Compatibility"
Cohesion: 0.25
Nodes (8): Architectures, BPF-LSM must be active, not just compiled in, Check your host first, Compatibility, Distributions, Kernel floors per mode, Low-end hardware, Mandatory and best-effort hooks

### Community 181 - "How it works"
Cohesion: 0.25
Nodes (8): eBPF at three levels, Fail closed, Global flags, How it works, Identity is the executable's inode, Output, Watching from a browser: `--serve`, Where it sits: DAC, MAC and app-listener

### Community 183 - "launch_rule_at"
Cohesion: 0.50
Nodes (3): launch_rule_at(), launch_step(), launch_token()

### Community 184 - "loadDaemonConfig"
Cohesion: 0.29
Nodes (7): loadDaemonConfig(), noteReloadedSettings(), resolveConfigPath(), TestLoadDaemonConfigEmptyResourcesIsCriticalStartup(), TestLoadDaemonConfigMissingFileIsCriticalStartup(), ConfigureInodeCapacity(), TestConfigureInodeCapacityBounds()

### Community 186 - "Binary trust"
Cohesion: 0.29
Nodes (7): Binary ledger, Binary replacement, Binary trust, Hash cache, Launch scan, Multicall binaries, trust-binaries

### Community 187 - "guard: block file access by program"
Cohesion: 0.33
Nodes (6): Executables inside the guarded tree, Flags, guard: block file access by program, Multicall binaries, Replaced binaries, Whitelist and blacklist

### Community 188 - "network-guard: block network access by program"
Cohesion: 0.33
Nodes (6): Code integrity (`-w`), Flags, Multicall binaries, network-guard: block network access by program, Replaced binaries, What code integrity cannot stop

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

### Community 196 - "batchInodeRows"
Cohesion: 0.19
Nodes (12): markReady(), notifySystemdReady(), inodeBatchFn, appendInodeRow(), batchInodeRows(), inodeRowsOf(), fakeBatches(), TestBatchInodeRows() (+4 more)

### Community 197 - "monitor: see what touches a path"
Cohesion: 0.67
Nodes (3): Events, Flags, monitor: see what touches a path

### Community 199 - "Execute"
Cohesion: 0.18
Nodes (10): installHooks(), releaseBuild(), Requested(), Run(), TestRequested(), TestSubcommandOutranksGUI(), Execute(), initLogger() (+2 more)

### Community 206 - "Desktop GUI"
Cohesion: 0.22
Nodes (9): Actions, Desktop GUI, Files, Home, Install, Network, Tables, The daemon binary (+1 more)

## Knowledge Gaps
- **32 isolated node(s):** `github.com/Virgula0/app-listener`, `IntegrationSuite`, `open_by_handle_at POC`, `statonly stat/statx metadata leak`, `app-listener-gui best-effort build` (+27 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 520 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **32 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `Process inspectors` connect `time.Time` to `Daemon configuration`?**
  _High betweenness centrality (0.031) - this node is a cross-community bridge._
- **What connects `github.com/Virgula0/app-listener`, `IntegrationSuite`, `open_by_handle_at POC` to the rest of the system?**
  _32 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `daemon/daemon.go` be split into smaller, more focused modules?**
  _Cohesion score 0.10435475670253261 - nodes in this community are weakly interconnected._
- **Why does `Daemon configuration` connect `Daemon configuration` to `time.Time`, `README.md`?**
  _High betweenness centrality (0.031) - this node is a cross-community bridge._
- **Should `testing.T` be split into smaller, more focused modules?**
  _Cohesion score 0.02182468694096601 - nodes in this community are weakly interconnected._
- **Why does `absPath()` connect `absPath` to `IntegrationSuite`, `IntegrationSuite`, `main_test.go`, `IntegrationSuite`, `IntegrationSuite`, `daemon_multicall_test.go`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `.netIntegrityContainer`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`?**
  _High betweenness centrality (0.031) - this node is a cross-community bridge._
- **Should `IntegrationSuite` be split into smaller, more focused modules?**
  _Cohesion score 0.08958130477117819 - nodes in this community are weakly interconnected._